package db_mock

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tomasdembelli/course-manager/models"
)

const (
	mockErrMessage = "mock error"
)

type mockError struct {
	message string
}

func NewMockError() *mockError {
	return &mockError{message: mockErrMessage}
}

func (e *mockError) Error() string {
	return e.message
}

func (e *mockError) Is(target error) bool {
	return target.Error() == e.message
}

type Config struct {
	CourseByUUID map[uuid.UUID]models.Course
	ErrByTutor   error
	ErrByStudent error
	ErrCreate    error
	ErrById      error
	ErrUpdate    error
	ErrDelete    error
	ErrList      error
}

type MockRepo struct {
	courseByUUID map[uuid.UUID]models.Course
	errByTutor   error
	errByStudent error
	errCreate    error
	errById      error
	errUpdate    error
	errDelete    error
	errList      error
}

func NewMockRepo(config *Config) *MockRepo {
	if config == nil {
		return &MockRepo{
			courseByUUID: make(map[uuid.UUID]models.Course),
		}
	}
	return &MockRepo{
		courseByUUID: config.CourseByUUID,
		errByTutor:   config.ErrByTutor,
		errByStudent: config.ErrByStudent,
		errCreate:    config.ErrCreate,
		errById:      config.ErrById,
		errUpdate:    config.ErrUpdate,
		errDelete:    config.ErrDelete,
		errList:      config.ErrList,
	}
}

func (m *MockRepo) safeInit() {
	if m.courseByUUID == nil {
		m.courseByUUID = make(map[uuid.UUID]models.Course)
	}
}

func (m *MockRepo) ById(_ context.Context, courseUuid uuid.UUID) (*models.Course, error) {
	if m.errById != nil {
		return nil, m.errById
	}
	course, found := m.courseByUUID[courseUuid]
	if !found {
		return nil, fmt.Errorf("course not found (UUID: %v)", courseUuid)
	}
	return &course, nil
}

func (m *MockRepo) ByTutor(_ context.Context, tutorUuid uuid.UUID) ([]models.Course, error) {
	if m.errByTutor != nil {
		return nil, m.errByTutor
	}
	var result []models.Course
	for _, course := range m.courseByUUID {
		if course.Tutor.Uuid == tutorUuid {
			result = append(result, course)
		}
	}
	return result, nil
}

func (m *MockRepo) ByStudent(_ context.Context, studentUuid uuid.UUID) ([]models.Course, error) {
	if m.errByStudent != nil {
		return nil, m.errByStudent
	}
	var result []models.Course
	for _, course := range m.courseByUUID {
		for _, student := range course.Students {
			if student.Uuid == studentUuid {
				result = append(result, course)
			}
		}
	}
	return result, nil
}

func (m *MockRepo) List(_ context.Context) ([]models.Course, error) {
	if m.errList != nil {
		return nil, m.errList
	}
	var courses = make([]models.Course, len(m.courseByUUID))
	var counter int
	for _, course := range m.courseByUUID {
		courses[counter] = course
		counter++
	}
	return courses, nil
}

func (m *MockRepo) Create(_ context.Context, course models.Course) error {
	m.safeInit()
	if m.errCreate != nil {
		return m.errCreate
	}
	m.courseByUUID[course.Uuid] = course
	return nil
}

func (m *MockRepo) Delete(_ context.Context, courseUUID uuid.UUID) error {
	m.safeInit()
	if m.errDelete != nil {
		return m.errDelete
	}
	delete(m.courseByUUID, courseUUID)
	return nil
}

func (m *MockRepo) Update(_ context.Context, course models.Course) error {
	m.safeInit()
	if m.errUpdate != nil {
		return NewMockError()
	}
	m.courseByUUID[course.Uuid] = course
	return nil
}
