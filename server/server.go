package server

import (
	"github.com/labstack/echo/v4"
	db_mock "github.com/tomasdembelli/course-manager/db-mock"
	"github.com/tomasdembelli/course-manager/services"
	"log"
	"os"
	"strconv"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

const devEnvironment = "development"

type Config struct {
	Port   int
	Logger *log.Logger
}

func StartServer(config *Config) {
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
	apiV1, err := NewApiV1(&courseManager, config.Logger)
	if err != nil {
		log.Fatal("unable to start apiV1", err)
	}

	e := echo.New()
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     echoMiddleware.DefaultCORSConfig.AllowMethods,
		AllowCredentials: true,
	}))
	v1 := e.Group("/v1")
	apiV1.Attach(v1)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
}
