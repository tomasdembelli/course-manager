package main

import (
	"github.com/tomasdembelli/course-manager/db-mock"
	"github.com/tomasdembelli/course-manager/server"
	"github.com/tomasdembelli/course-manager/services"
	"log"
)

func main() {
	mockRepo := db_mock.NewMockRepo(&db_mock.Config{
		CourseByUUID: db_mock.CourseByUUID,
	})
	courseManager, err := services.NewCourseManager(mockRepo, nil)
	if err != nil {
		log.Fatalf("unable to start course manager service")
	}
	server.StartServer(8000, &courseManager)
}
