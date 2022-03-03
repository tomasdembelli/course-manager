package models

// Student defines a student who registers a course.
// It embeds the User type.
type Student struct {
	User
	Faculty notEmptyString `json:"faculty"`
}
