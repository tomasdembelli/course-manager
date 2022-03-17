package models

import (
	"fmt"

	"github.com/google/uuid"
)

const canNotBeEmptyFmt = "%s cannot be empty"

// User constitutes a base model for the person who interacts with a course.
type User struct {
	Uuid     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	Lastname string    `json:"lastname"`
}

// Valid returns an error when the name and/or lastname are empty strings.
func (u *User) Valid() error {
	if u.Name == "" {
		return fmt.Errorf(canNotBeEmptyFmt, "name")
	}
	if u.Lastname == "" {
		return fmt.Errorf(canNotBeEmptyFmt, "lastname")
	}
	return nil
}
