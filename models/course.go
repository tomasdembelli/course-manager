package models

import "github.com/google/uuid"

type CourseMeta struct {
	Uuid  uuid.UUID `json:"uuid,omitempty"`
	Name  string    `json:"name"`
	Tutor *Tutor    `json:"tutor"`
}

// Course defines a course.
type Course struct {
	CourseMeta
	Students map[uuid.UUID]Student `json:"students"`
}
