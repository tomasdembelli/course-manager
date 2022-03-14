package db_mock

import (
	"github.com/google/uuid"
	"github.com/tomasdembelli/course-manager/models"
)

var CourseByUUID = map[uuid.UUID]models.Course{
	uuid.New(): {
		Uuid: uuid.New(),
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
