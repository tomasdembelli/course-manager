package models

// Course defines a course.
type Course struct {
	Name     notEmptyString `json:"name"`
	Tutor    Tutor          `json:"tutor"`
	Students []Student      `json:"students"`
}
