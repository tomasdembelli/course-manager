package server

import (
	"github.com/labstack/echo/v4"
	"github.com/tomasdembelli/course-manager/services"
	"log"
	"strconv"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

type Config struct {
	Port             int
	CourseManagerSvc *services.CourseManager
}

func StartServer(config *Config) {
	apiV1, err := NewApiV1(config.CourseManagerSvc)
	if err != nil {
		log.Fatal("unable to start apiV1", err)
	}

	e := echo.New()
	e.Use(echoMiddleware.Logger())
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     echoMiddleware.DefaultCORSConfig.AllowMethods,
		AllowCredentials: true,
	}))
	v1 := e.Group("/v1")
	apiV1.Attach(v1)

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(config.Port)))
}
