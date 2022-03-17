package models

import "github.com/google/uuid"

// Course defines a course.
type Course struct {
	Uuid     uuid.UUID             `json:"uuid"`
	Name     string                `json:"name"`
	Tutor    *Tutor                `json:"tutor"`
	Students map[uuid.UUID]Student `json:"students"`
}
