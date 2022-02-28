package models

type Student struct {
	User
	Faculty notEmptyString `json:"faculty"`
}
