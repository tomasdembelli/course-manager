package models

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var canNotBeEmptyErr = errors.New("field can not be empty")

// notEmptyString is a wrapper around the string type.
type notEmptyString string

// User constitutes a base model for the person who interacts with a course.
type User struct {
	Uuid     uuid.UUID      `json:"uuid"`
	Name     notEmptyString `json:"name"`
	Lastname notEmptyString `json:"lastname"`
}

// UnmarshalJSON is called when a string of byte is unmarsalled into a notEmptyString type variable.
// It returns an error when the given data is an empty string.
func (s *notEmptyString) UnmarshalJSON(data []byte) error {
	var ns string
	if err := json.Unmarshal(data, &ns); err != nil {
		return err
	}
	if ns == "" {
		return canNotBeEmptyErr
	}
	return nil
}
