package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestStudent_MarshalJSON(t *testing.T) {
	student := Student{
		User: User{
			Uuid:     uuid.MustParse("c46358be-a216-4083-8bc2-0c4eda703b4a"),
			Name:     "John",
			Lastname: "Doe",
		},
		Faculty: "Electronic",
	}
	got, err := json.Marshal(student)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	want := []byte(`{"uuid":"c46358be-a216-4083-8bc2-0c4eda703b4a","name":"John","lastname":"Doe","faculty":"Electronic"}`)
	if string(got) != string(want) {
		t.Errorf("expected %v, got %v", string(want), string(got))
	}
}

func TestTutor_MarshalJSON(t *testing.T) {
	tutor := Tutor{
		User: User{
			Uuid:     uuid.MustParse("c46358be-a216-4083-8bc2-0c4eda703b4a"),
			Name:     "John",
			Lastname: "Doe",
		},
		Faculty:    "Electronic",
		LecturerOf: "Applied math",
	}
	got, err := json.Marshal(tutor)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	want := []byte(`{"uuid":"c46358be-a216-4083-8bc2-0c4eda703b4a","name":"John","lastname":"Doe","faculty":"Electronic","lecturerOf":"Applied math"}`)
	if string(got) != string(want) {
		t.Errorf("expected %v, got %v", string(want), string(got))
	}
}
