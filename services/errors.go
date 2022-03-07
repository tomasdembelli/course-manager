package services

import (
	"fmt"
	"github.com/google/uuid"
)

const (
	courseNotFoundFmt = "Course with UUID = %d not found"
	cannotBeNilFmt    = "%v cannot be nil"
	validationErrFmt  = "validation failed: %v"
)

// NotFoundError should be returned when a service can't find a course.
type NotFoundError struct {
	message string
}

func NewNotFoundErr(message string) *NotFoundError {
	return &NotFoundError{message: message}
}

func NewOrderNotFoundErr(courseUuid uuid.UUID) *NotFoundError {
	return NewNotFoundErr(fmt.Sprintf(courseNotFoundFmt, courseUuid))
}

// Error implements error. Returns the error message associated with the NilErr.
func (e *NotFoundError) Error() string {
	return e.message
}

// Is reports whether the given error is equal to the NilErr
func (e *NotFoundError) Is(target error) bool { return target.Error() == e.message }

// NilErr should be returned when an input is nil.
type NilErr struct {
	message string
}

func NewNilErr(item string) *NilErr {
	return &NilErr{message: fmt.Sprintf(cannotBeNilFmt, item)}
}

// Error implements error. Returns the error message associated with the NilErr.
func (e *NilErr) Error() string {
	return e.message
}

// Is reports whether the given error is equal to the NilErr
func (e *NilErr) Is(target error) bool { return target.Error() == e.message }

type courseConstraint string

const (
	tutorMaxCourse   courseConstraint = "a tutor can facilitate maximum 2 courses"
	studentMaxCourse courseConstraint = "a student can register to maximum 4 courses"
	courseMaxStudent courseConstraint = "maximum 20 students can register a course"
)

type CourseConstraintErr struct {
	message string
}

func NewCourseConstraintErr(check courseConstraint) *CourseConstraintErr {
	return &CourseConstraintErr{message: fmt.Sprintf(validationErrFmt, check)}
}

func (e *CourseConstraintErr) Error() string {
	return e.message
}

func (e *CourseConstraintErr) Is(target error) bool {
	return target.Error() == e.message
}
