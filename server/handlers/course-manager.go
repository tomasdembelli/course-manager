package handlers

import (
	"fmt"
	"github.com/tomasdembelli/course-manager/services"
)

// Api exposes a services.CourseManager via HTTP endpoints.
type Api struct {
	courseManagerSvc *services.CourseManager
}

// NewApi returns a new API that wraps the given services.CourseManager with HTTP endpoints.
func NewApi(courseManager *services.CourseManager) (*Api, error) {
	if courseManager == nil {
		return nil, fmt.Errorf("coursse manager cannot be nil")
	}
	return &Api{
		courseManagerSvc: courseManager,
	}, nil
}
