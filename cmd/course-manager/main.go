package main

import (
	db_mock "github.com/tomasdembelli/course-manager/db-mock"
	"github.com/tomasdembelli/course-manager/server"
	"github.com/tomasdembelli/course-manager/services"
	"log"
	"os"
)

const devEnvironment = "development"

func main() {
	var repo services.Repo
	if os.Getenv("ENVIRONMENT") == devEnvironment {
		repo = db_mock.NewMockRepo(&db_mock.Config{
			CourseByUUID: db_mock.CourseByUUID,
		})
	}
	courseManager, err := services.NewCourseManager(repo, log.Default())
	if err != nil {
		log.Fatalf("unable to start course manager service %v", err)
	}
	server.StartServer(&server.Config{
		Port:             8000,
		CourseManagerSvc: &courseManager,
	})
}
