package smoke_tests

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (r *RequestParams) Do() error {
	data, err := json.Marshal(&r.Payload)
	if err != nil {
		return fmt.Errorf("failed to marshall payload %v", r.Payload)
	}
	request, err := http.NewRequest(strings.ToUpper(r.Method), r.BaseUrl+r.Path, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create a %v request to %v", request.Method, request.RequestURI)
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
		return fmt.Errorf("failed to DO a %v request to %v", request.Method, request.RequestURI)
	}
	r.StatusCode = response.StatusCode
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response %v", err)
	}
	if string(body) != "" {
		err = json.Unmarshal(body, &r.ResponseBody)
		if err != nil {
			return fmt.Errorf("failed to unmarshall response responseBody %v, %v", r.Path, err)
		}
	}
	return nil
}
