package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	db_mock "github.com/tomasdembelli/course-manager/db-mock"
	"github.com/tomasdembelli/course-manager/services"
)

var (
	courseManager services.CourseManager
)

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	res := &events.APIGatewayProxyResponse{}
	switch req.Path {
	case "/v1/listCourses":
		courses, err := courseManager.List(ctx)
		if err != nil {
			return res, err
		}

		body, err := json.Marshal(courses)
		if err != nil {
			return res, err
		}

		res.Body = string(body)
		return res, nil

	default:
		return res, fmt.Errorf("unrecognized path %s", req.Path)

	}
}

func main() {
	//TODO: fix this
	repo := db_mock.NewMockRepo(&db_mock.Config{
		CourseByUUID: db_mock.CourseByUUID,
	})

	var err error
	courseManager, err = services.NewCourseManager(repo, log.Default())
	if err != nil {
		log.Fatalf("unable to start course manager service %v", err)
	}

	lambda.Start(handler)

}
