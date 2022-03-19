package models

import (
	"github.com/google/uuid"
)

const canNotBeEmptyFmt = "%s cannot be empty"

// User constitutes a base model for the person who interacts with a course.
type User struct {
	Uuid     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
}
