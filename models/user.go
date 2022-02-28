package models

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

var canNotBeEmptyErr = errors.New("field can not be empty")

type notEmptyString string

// User constitutes a base model for the person who interacts with a course.
type User struct {
	Uuid     uuid.UUID      `json:"uuid"`
	Name     notEmptyString `json:"name"`
	Lastname notEmptyString `json:"lastname"`
}

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
