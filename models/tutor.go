package models

// Tutor defines the attributes of a tutor who facilitates a course.
// It embeds the User type.
type Tutor struct {
	User
	Faculty    string `json:"faculty"`
	LecturerOf string `json:"lecturerOf"`
}
