package db_mock

import (
	"github.com/google/uuid"
	"github.com/tomasdembelli/course-manager/models"
)

var CourseByUUID = map[uuid.UUID]models.Course{
	uuid.MustParse("2d2e10a1-94e2-4dff-a244-8733bee8b7a9"): {
		Uuid: uuid.MustParse("2d2e10a1-94e2-4dff-a244-8733bee8b7a9"),
		Name: "Mock Course",
		Tutor: &models.Tutor{
			User: models.User{
				Uuid:     uuid.New(),
				Name:     "Mock",
				Lastname: "Tutor",
			},
			Faculty:    "Mock Faculty",
			LecturerOf: "Mock Lecturer of",
		},
		Students: map[uuid.UUID]models.Student{
			uuid.New(): {
				User: models.User{
					Uuid:     uuid.New(),
					Name:     "Mock",
					Lastname: "Student",
				},
				Faculty: "Mock Faculty",
			},
		},
	},
}
