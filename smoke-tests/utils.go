package smoke_tests

import (
	"bytes"
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"strings"
)

type RequestParams struct {
	BaseUrl      string
	Path         string
	Header       map[string]string
	Method       string
	StatusCode   int
	ResponseBody interface{}
	Payload      map[string]interface{}
	Params       map[string]string
}

func (r *RequestParams) Do() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	data, err := json.Marshal(&r.Payload)
	if err != nil {
		logger.Fatal("failed to marshall payload %v", zap.Any("payload", r.Payload))
	}
	request, err := http.NewRequest(strings.ToUpper(r.Method), r.BaseUrl+r.Path, bytes.NewBuffer(data))
	if err != nil {
		logger.Fatal("failed to create a %v request to %v", zap.String("method", request.Method), zap.String("uri", request.RequestURI))
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	request.Header.Set("Accept", "application/json; charset=UTF-8")
	if r.Header != nil {
		for k, v := range r.Header {
			request.Header.Set(k, v)
		}
	}
	if r.Params != nil {
		q := request.URL.Query()
		for k, v := range r.Params {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		logger.Fatal("failed to DO a %v request to %v", zap.String("method", request.Method), zap.String("uri", request.RequestURI))
	}
	r.StatusCode = response.StatusCode
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		logger.Fatal("error reading response %v", zap.Error(err))
	}
	if string(body) != "" {
		err = json.Unmarshal(body, &r.ResponseBody)
		if err != nil {
			logger.Warn("failed to unmarshall response responseBody", zap.String("path", r.Path), zap.Error(err))
		}
	}
}
