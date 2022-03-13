package server

import (
	"github.com/google/uuid"
	"github.com/tomasdembelli/course-manager/models"
)

// Course defines a course object used in HTTP server.
type Course struct {
	Name     string
	Tutor    models.Tutor
	Students []models.Student
}

// CreateCourse should be used at the HTTP endpoint for creating a course.
type CreateCourse struct {
	Course Course `form:"course"`
}

// CourseByUUID should be used at the HTTP endpoint for querying an individual course by its UUID.
type CourseByUUID struct {
	UUID uuid.UUID `param:"courseUUID"`
}

// RegisterStudent should be used at the HTTP endpoint registering a student to a given course.
type RegisterStudent struct {
	CourseUUID uuid.UUID      `param:"courseUUID"`
	Student    models.Student `form:"student"`
}

// UnregisterStudent should be used at the HTTP endpoint unregistering a student from a given course.
type UnregisterStudent struct {
	CourseUUID  uuid.UUID `param:"courseUUID"`
	StudentUUID uuid.UUID `form:"studentUUID"`
}
