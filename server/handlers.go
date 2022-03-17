package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tomasdembelli/course-manager/services"
	"log"
	"net/http"
)

var (
	notFoundMessage = map[string]string{"message": "not found"}
	successMessage  = map[string]string{"message": "successful"}
)

// ApiV1 exposes a services.CourseManager via HTTP endpoints.
type ApiV1 struct {
	courseManagerSvc *services.CourseManager
	logger           *log.Logger
}

// NewApiV1 returns a new API that wraps the given services.CourseManager with HTTP endpoints.
func NewApiV1(courseManager *services.CourseManager, logger *log.Logger) (*ApiV1, error) {
	if courseManager == nil {
		return nil, fmt.Errorf("coursse manager cannot be nil")
	}
	if logger == nil {
		logger = log.Default()
	}
	return &ApiV1{
		courseManagerSvc: courseManager,
		logger:           logger,
	}, nil
}

func (a *ApiV1) Attach(group *echo.Group) {
	group.GET("/listCourses", a.ListCourses)
	group.GET("/getCourse/:courseUUID", a.GetCourse)
	group.PUT("/registerStudent/:courseUUID", a.RegisterStudent)
	group.PUT("/unregisterStudent/:courseUUID", a.UnregisterStudent)
	group.POST("/createCourse", a.Create)
}

func (a *ApiV1) ListCourses(ec echo.Context) error {
	courses, err := a.courseManagerSvc.List(ec.Request().Context())
	if err != nil {
		return err
	}
	return ec.JSON(http.StatusOK, courses)
}

func (a *ApiV1) GetCourse(ec echo.Context) error {
	request := new(CourseByUUID)
	if err := ec.Bind(request); err != nil {
		a.logger.Printf("unable to bind request params: %v", err)
		return err
	}
	course, err := a.courseManagerSvc.Get(ec.Request().Context(), request.UUID)
	if err != nil {
		a.logger.Printf("unable to retrieve the course: %v", err)
		return ec.JSON(http.StatusNotFound, notFoundMessage)
	}
	return ec.JSON(http.StatusOK, course)
}

func (a *ApiV1) RegisterStudent(ec echo.Context) error {
	request := new(RegisterStudent)
	if err := ec.Bind(request); err != nil {
		a.logger.Printf("unable to bind request params: %v", err)
		return err
	}
	err := a.courseManagerSvc.RegisterStudent(ec.Request().Context(), request.CourseUUID, request.Student)
	if err != nil {
		a.logger.Println(err)
		return ec.JSON(http.StatusBadRequest, map[string]string{"message": "unable to register student"})
	}
	return ec.JSON(http.StatusOK, successMessage)
}

func (a *ApiV1) UnregisterStudent(ec echo.Context) error {
	request := new(UnregisterStudent)
	if err := ec.Bind(request); err != nil {
		a.logger.Printf("unable to bind request params: %v", err)
		return err
	}
	err := a.courseManagerSvc.UnregisterStudent(ec.Request().Context(), request.CourseUUID, request.StudentUUID)
	if err != nil {
		a.logger.Println(err)
		return ec.JSON(http.StatusBadRequest, map[string]string{"message": "unable to unregister student"})
	}
	return ec.JSON(http.StatusOK, successMessage)
}

func (a *ApiV1) Create(ec echo.Context) error {
	request := new(CreateCourse)
	if err := ec.Bind(request); err != nil {
		a.logger.Printf("unable to bind request params: %v", err)
		return err
	}
	course, err := a.courseManagerSvc.Create(ec.Request().Context(), request.Course)
	if err != nil {
		a.logger.Println(err)
		return ec.JSON(http.StatusBadRequest, map[string]string{"message": "unable to create the course"})
	}
	return ec.JSON(http.StatusCreated, course)
}
