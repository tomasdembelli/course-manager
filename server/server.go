package server

import (
	"github.com/labstack/echo/v4"
	"github.com/tomasdembelli/course-manager/services"
	"net/http"
	"strconv"

	echoMiddleware "github.com/labstack/echo/v4/middleware"
)

func StartServer(port int, courseManager *services.CourseManager) {
	_ = courseManager
	e := echo.New()
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     echoMiddleware.DefaultCORSConfig.AllowMethods,
		AllowCredentials: true,
	}))
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
}
