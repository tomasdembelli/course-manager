package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/tomasdembelli/course-manager/models"
)

const (
	mockErrMessage = "mock error"
)

var fixedUuid = uuid.MustParse("5d61cbc8-9ccd-4348-a623-d61dd7658dd7")

type mockError struct {
	message string
}

func newMockError() *mockError {
	return &mockError{message: mockErrMessage}
}

func (e *mockError) Error() string {
	return e.message
}

func (e *mockError) Is(target error) bool {
	return target.Error() == e.message
}

type mockRepo struct {
	courseByUUID map[uuid.UUID]models.Course
	errByTutor   error
	errByStudent error
	errCreate    error
	errById      error
	errUpdate    error
	errDelete    error
}

func (m *mockRepo) safeInit() {
	if m.courseByUUID == nil {
		m.courseByUUID = make(map[uuid.UUID]models.Course)
	}
}

func (m *mockRepo) ById(_ context.Context, courseUuid uuid.UUID) (models.Course, error) {
	if m.errById != nil {
		return models.Course{}, m.errById
	}
	return m.courseByUUID[courseUuid], nil
}

func (m *mockRepo) ByTutor(_ context.Context, tutorUuid uuid.UUID) ([]models.Course, error) {
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

func (m *mockRepo) ByStudent(_ context.Context, studentUuid uuid.UUID) ([]models.Course, error) {
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

func (m *mockRepo) List(_ context.Context) ([]models.Course, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) Create(_ context.Context, course models.Course) error {
	m.safeInit()
	if m.errCreate != nil {
		return m.errCreate
	}
	m.courseByUUID[course.Uuid] = course
	return nil
}

func (m *mockRepo) Delete(_ context.Context, courseUUID uuid.UUID) error {
	m.safeInit()
	if m.errDelete != nil {
		return m.errDelete
	}
	delete(m.courseByUUID, courseUUID)
	return nil
}

func (m *mockRepo) Update(_ context.Context, course models.Course) error {
	m.safeInit()
	if m.errUpdate != nil {
		return newMockError()
	}
	m.courseByUUID[course.Uuid] = course
	return nil
}

func TestNewCourseManager(t *testing.T) {
	type args struct {
		repo   Repo
		logger *log.Logger
	}
	tests := []struct {
		name        string
		args        args
		want        CourseManager
		wantErr     bool
		expectedErr error
	}{
		{
			name: "course manager with nil repo and nil logger should error",
			args: args{
				repo:   nil,
				logger: nil,
			},
			want:        CourseManager{},
			wantErr:     true,
			expectedErr: NewNilErr("repo"),
		},
		{
			name: "course manager with a compliant repo and nil logger should pass",
			args: args{
				repo:   &mockRepo{},
				logger: nil,
			},
			want: CourseManager{
				repo:   &mockRepo{},
				logger: log.Default(),
			},
		},
		{
			name: "course manager with a compliant repo and default logger should pass",
			args: args{
				repo:   &mockRepo{},
				logger: log.Default(),
			},
			want: CourseManager{
				repo:   &mockRepo{},
				logger: log.Default(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCourseManager(tt.args.repo, tt.args.logger)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCourseManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCourseManager() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr && !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("NewCourseManager() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func generateUsersInCourse(numberOfStudents int) models.Course {
	course := models.Course{
		Name: "test course",
		Uuid: fixedUuid,
		Tutor: &models.Tutor{
			User: models.User{
				Uuid: fixedUuid,
			},
		},
		Students: make(map[uuid.UUID]models.Student),
	}

	for i := 0; i < numberOfStudents; i++ {
		studentUuid := uuid.New()
		course.Students[studentUuid] = models.Student{
			User: models.User{
				Uuid: studentUuid,
			},
		}
	}
	return course
}

func TestCourseManager_Create(t *testing.T) {
	predefinedCourse := generateUsersInCourse(10)
	var predefinedStudents []uuid.UUID
	for studentUuid := range predefinedCourse.Students {
		predefinedStudents = append(predefinedStudents, studentUuid)
	}
	delete(predefinedCourse.Students, predefinedStudents[0])
	refinedStudents := predefinedCourse.Students

	type fields struct {
		repo   Repo
		logger *log.Logger
	}
	type args struct {
		ctx    context.Context
		course models.Course
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        *models.Course
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "nil tutor error",
			fields:      fields{repo: &mockRepo{}},
			args:        args{ctx: context.TODO(), course: models.Course{}},
			want:        nil,
			wantErr:     true,
			expectedErr: NewNilErr("tutor"),
		},
		{
			name: "error at ByTutor",
			fields: fields{repo: &mockRepo{
				errByTutor: newMockError(),
			}},
			args: args{ctx: context.TODO(), course: models.Course{
				Tutor: &models.Tutor{
					User: models.User{
						Uuid: fixedUuid,
					},
				},
			}},
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to retrieve courses: %w", newMockError()),
		},
		{
			name: "tutor max course validation",
			fields: fields{repo: &mockRepo{
				courseByUUID: map[uuid.UUID]models.Course{
					uuid.New(): {
						Tutor: &models.Tutor{
							User: models.User{Uuid: fixedUuid},
						},
						Name: "course 1",
					},
					uuid.New(): {
						Tutor: &models.Tutor{
							User: models.User{Uuid: fixedUuid},
						},
						Name: "course 2",
					},
				},
			}},
			args: args{ctx: context.TODO(), course: models.Course{
				Tutor: &models.Tutor{
					User: models.User{
						Uuid: fixedUuid,
					},
				},
			}},
			want:        nil,
			wantErr:     true,
			expectedErr: NewCourseConstraintErr(tutorMaxCourse),
		},
		{
			name: "course max student validation",
			fields: fields{
				repo: &mockRepo{},
			},
			args: args{
				ctx:    context.TODO(),
				course: generateUsersInCourse(21),
			},
			want:        nil,
			wantErr:     true,
			expectedErr: NewCourseConstraintErr(courseMaxStudent),
		},
		{
			name: "err at ByStudent",
			fields: fields{
				repo: &mockRepo{
					errByStudent: newMockError(),
				},
			},
			args: args{
				ctx:    context.TODO(),
				course: generateUsersInCourse(10),
			},
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to retrieve courses: %w", newMockError()),
		},
		{
			name: "1 ineligible students out of 10",
			fields: fields{
				repo: &mockRepo{
					courseByUUID: map[uuid.UUID]models.Course{
						uuid.New(): {
							Name: "course 1",
							Tutor: &models.Tutor{
								User: models.User{
									Uuid: uuid.New(),
								},
							},
							Students: map[uuid.UUID]models.Student{
								predefinedStudents[0]: {},
							},
						},
						uuid.New(): {
							Name: "course 2",
							Tutor: &models.Tutor{
								User: models.User{
									Uuid: uuid.New(),
								},
							},
							Students: map[uuid.UUID]models.Student{
								predefinedStudents[0]: {},
							},
						},
						uuid.New(): {
							Name: "course 3",
							Tutor: &models.Tutor{
								User: models.User{
									Uuid: uuid.New(),
								},
							},
							Students: map[uuid.UUID]models.Student{
								predefinedStudents[0]: {},
							},
						},
						uuid.New(): {
							Name: "course 4",
							Tutor: &models.Tutor{
								User: models.User{
									Uuid: uuid.New(),
								},
							},
							Students: map[uuid.UUID]models.Student{
								predefinedStudents[0]: {},
							},
						},
					},
				},
			},
			args: args{
				ctx:    context.TODO(),
				course: predefinedCourse,
			},
			want: &models.Course{
				Uuid: fixedUuid,
				Name: "test course",
				Tutor: &models.Tutor{
					User: models.User{
						Uuid: fixedUuid,
					},
				},
				Students: refinedStudents,
			},
		},
		{
			name: "err at Create",
			fields: fields{
				repo: &mockRepo{
					errCreate: newMockError(),
				},
			},
			args: args{
				ctx:    context.TODO(),
				course: predefinedCourse,
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to create the course: %w", newMockError()),
		},
		{
			name: "err at ById",
			fields: fields{
				repo: &mockRepo{
					errById: newMockError(),
				},
			},
			args: args{
				ctx:    context.TODO(),
				course: predefinedCourse,
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to retrieve the course: %w", newMockError()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseManager, err := NewCourseManager(tt.fields.repo, tt.fields.logger)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := courseManager.Create(tt.args.ctx, tt.args.course)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
			if tt.wantErr && !reflect.DeepEqual(err, tt.expectedErr) {
				t.Errorf("Create() error = %v, expectedErr = %v", err, tt.expectedErr)
			}
		})
	}
}

func TestCourseManager_RegisterStudent(t *testing.T) {
	predefinedCourse := generateUsersInCourse(10)
	type fields struct {
		repo   Repo
		logger *log.Logger
	}
	type args struct {
		ctx        context.Context
		courseUUID uuid.UUID
		student    models.Student
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantErr            bool
		expectedErrMessage string
	}{
		{
			name: "error at ById",
			fields: fields{
				repo: &mockRepo{
					errById: newMockError(),
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: uuid.New(),
				student:    models.Student{},
			},
			wantErr:            true,
			expectedErrMessage: "unable to retrieve the course: mock error",
		},
		{
			name: "error at Update",
			fields: fields{
				repo: &mockRepo{
					courseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourse,
					},
					errUpdate: newMockError(),
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: fixedUuid,
				student: models.Student{
					User: models.User{Uuid: fixedUuid},
				},
			},
			wantErr:            true,
			expectedErrMessage: "unable to update the course: mock error",
		},
		{
			name: "successful registry",
			fields: fields{
				repo: &mockRepo{
					courseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourse,
					},
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: fixedUuid,
				student: models.Student{
					User: models.User{Uuid: fixedUuid},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCourseManager(tt.fields.repo, tt.fields.logger)
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			err = c.RegisterStudent(tt.args.ctx, tt.args.courseUUID, tt.args.student)
			if tt.wantErr && tt.expectedErrMessage != err.Error() {
				t.Errorf("RegisterStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("unexpected error RegisterStudent() error = %v", err)
				}
				courses, err := c.repo.ByStudent(tt.args.ctx, tt.args.student.Uuid)
				if err != nil {
					t.Fatal("unexpected error", err)
				}
				if courses[0].Students[tt.args.student.Uuid].Uuid != tt.args.student.Uuid {
					t.Errorf("failed to register the student. expected %v, got %v", tt.args.student.Uuid, courses[0].Students[tt.args.student.Uuid].Uuid)
				}
			}
		})
	}
}

func TestCourseManager_UnregisterStudent(t *testing.T) {
	predefinedCourse := generateUsersInCourse(10)
	var anExistingStudent models.Student
	for _, student := range predefinedCourse.Students {
		anExistingStudent = student
		break
	}
	type fields struct {
		repo   Repo
		logger *log.Logger
	}
	type args struct {
		ctx        context.Context
		courseUUID uuid.UUID
		student    models.Student
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantErr            bool
		expectedErrMessage string
	}{
		{
			name: "error at ById",
			fields: fields{
				repo: &mockRepo{
					errById: newMockError(),
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: uuid.New(),
				student:    models.Student{},
			},
			wantErr:            true,
			expectedErrMessage: "unable to retrieve the course: mock error",
		},
		{
			name: "error at Update",
			fields: fields{
				repo: &mockRepo{
					courseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourse,
					},
					errUpdate: newMockError(),
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: fixedUuid,
				student: models.Student{
					User: models.User{Uuid: fixedUuid},
				},
			},
			wantErr:            true,
			expectedErrMessage: "unable to update the course: mock error",
		},
		{
			name: "successful un-registry",
			fields: fields{
				repo: &mockRepo{
					courseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourse,
					},
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: fixedUuid,
				student:    anExistingStudent,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCourseManager(tt.fields.repo, tt.fields.logger)
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			err = c.UnregisterStudent(tt.args.ctx, tt.args.courseUUID, tt.args.student)
			if tt.wantErr && tt.expectedErrMessage != err.Error() {
				t.Errorf("UnregisterStudent() error = %v, wantErr %v", err.Error(), tt.expectedErrMessage)
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("unexpected error UnregisterStudent() error = %v", err)
				}
				courses, err := c.repo.ByStudent(tt.args.ctx, tt.args.student.Uuid)
				if err != nil {
					t.Fatal("unexpected error", err)
				}
				for _, course := range courses {
					if course.Uuid == tt.args.courseUUID {
						t.Errorf("failed to un-register the student from the course %v", course)
					}
				}

			}
		})
	}
}

func TestCourseManager_Delete(t *testing.T) {
	type fields struct {
		repo   Repo
		logger *log.Logger
	}
	type args struct {
		ctx        context.Context
		courseUUID uuid.UUID
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		wantErr            bool
		expectedErrMessage string
	}{
		{
			name: "error at repo Delete",
			fields: fields{
				repo: &mockRepo{
					errDelete: newMockError(),
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: uuid.New(),
			},
			wantErr:            true,
			expectedErrMessage: "unable to delete the course: mock error",
		},
		{
			name: "successful Delete",
			fields: fields{
				repo: &mockRepo{
					courseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: {},
					},
				},
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: fixedUuid,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCourseManager(tt.fields.repo, tt.fields.logger)
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			err = c.Delete(tt.args.ctx, tt.args.courseUUID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but none raised")
				}
				if tt.expectedErrMessage != err.Error() {
					t.Errorf("RegisterStudent() error = %v, wantErr %v", err.Error(), tt.expectedErrMessage)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error Delete() error = %v", err)
				}
				course, err := c.repo.ById(tt.args.ctx, tt.args.courseUUID)
				if err != nil {
					t.Fatal("unexpected error", err)
				}
				if course.Uuid == tt.args.courseUUID {
					t.Errorf("failed to delete course %v", tt.args.courseUUID)
				}

			}

		})
	}
}
