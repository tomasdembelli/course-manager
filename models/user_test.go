package models

import (
	"encoding/json"
	"testing"

	"github.com/google/uuid"
)

func TestUser_UnmarshalJSON(t *testing.T) {
	tests := map[string]struct {
		data        []byte
		shouldErr   bool
		expectedErr error
	}{
		"valid user": {
			data: []byte(`  {
				"uuid": "c46358be-a216-4083-8bc2-0c4eda703b4a",
				"name": "John",
				"lastname": "Doe"
			  }`,
			),
		},
		"invalid uuid": {
			data: []byte(`  {
				"uuid": "99999-a216-4083-8bc2-0c4eda703b4a",
				"name": "John",
				"lastname": "Doe"
			  }`,
			),
			shouldErr: true,
		},
		"empty uuid": {
			data: []byte(`  {
				"uuid": "",
				"name": "John",
				"lastname": "Doe"
			  }`,
			),
			shouldErr: true,
		},
		"empty name": {
			data: []byte(`  {
				"uuid": "c46358be-a216-4083-8bc2-0c4eda703b4a",
				"name": "",
				"lastname": "Doe"
			  }`,
			),
			shouldErr:   true,
			expectedErr: canNotBeEmptyErr,
		},
	}
	for desc, tt := range tests {
		t.Run(desc, func(t *testing.T) {
			var s User
			err := json.Unmarshal(tt.data, &s)
			if tt.shouldErr && err == nil {
				t.Errorf("unexpected error %v", err)
			}
			if tt.expectedErr != nil && tt.expectedErr != err {
				t.Errorf("expected %v, but got %v", tt.expectedErr, err)
			}
		})
	}
}

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
