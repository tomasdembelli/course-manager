package models

type Tutor struct {
	User
	Faculty    notEmptyString `json:"faculty"`
	LecturerOf notEmptyString `json:"lecturerOf"`
}
