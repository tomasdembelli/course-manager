package services

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/google/uuid"
	. "github.com/tomasdembelli/course-manager/db-mock"
	"github.com/tomasdembelli/course-manager/models"
)

var fixedUuid = uuid.MustParse("5d61cbc8-9ccd-4348-a623-d61dd7658dd7")

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
			name: "courseMeta manager with nil repo and nil logger should error",
			args: args{
				repo:   nil,
				logger: nil,
			},
			want:        CourseManager{},
			wantErr:     true,
			expectedErr: NewNilErr("repo"),
		},
		{
			name: "courseMeta manager with a compliant repo and nil logger should pass",
			args: args{
				repo:   &MockRepo{},
				logger: nil,
			},
			want: CourseManager{
				repo:   &MockRepo{},
				logger: log.Default(),
			},
		},
		{
			name: "courseMeta manager with a compliant repo and default logger should pass",
			args: args{
				repo:   &MockRepo{},
				logger: log.Default(),
			},
			want: CourseManager{
				repo:   &MockRepo{},
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
		CourseMeta: models.CourseMeta{
			Name: "test courseMeta",
			Uuid: fixedUuid,
			Tutor: &models.Tutor{
				User: models.User{
					Uuid: fixedUuid,
				},
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

	type fields struct {
		repo   Repo
		logger *log.Logger
	}
	type args struct {
		ctx        context.Context
		courseMeta models.CourseMeta
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
			fields:      fields{repo: &MockRepo{}},
			args:        args{ctx: context.TODO(), courseMeta: models.CourseMeta{}},
			want:        nil,
			wantErr:     true,
			expectedErr: NewNilErr("tutor"),
		},
		{
			name:   "error at ByTutor",
			fields: fields{repo: NewMockRepo(&Config{ErrByTutor: NewMockError()})},
			args: args{ctx: context.TODO(), courseMeta: models.CourseMeta{
				Tutor: &models.Tutor{
					User: models.User{
						Uuid: fixedUuid,
					},
				},
			},
			},
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to retrieve courses: %w", NewMockError()),
		},
		{
			name: "tutor max courseMeta validation",
			fields: fields{repo: NewMockRepo(&Config{
				CourseByUUID: map[uuid.UUID]models.Course{
					uuid.New(): {
						CourseMeta: models.CourseMeta{
							Tutor: &models.Tutor{
								User: models.User{Uuid: fixedUuid},
							},
							Name: "courseMeta 1",
						},
					},
					uuid.New(): {
						CourseMeta: models.CourseMeta{
							Tutor: &models.Tutor{
								User: models.User{Uuid: fixedUuid},
							},
							Name: "courseMeta 2",
						},
					},
				},
			})},
			args: args{ctx: context.TODO(), courseMeta: models.CourseMeta{
				Tutor: &models.Tutor{
					User: models.User{Uuid: fixedUuid},
				},
			},
			},
			want:        nil,
			wantErr:     true,
			expectedErr: NewCourseConstraintErr(tutorMaxCourseMsg),
		},
		{
			name: "err at Create",
			fields: fields{
				repo: NewMockRepo(&Config{
					ErrCreate: NewMockError(),
				}),
			},
			args: args{
				ctx:        context.TODO(),
				courseMeta: predefinedCourse.CourseMeta,
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to create the course: %w", NewMockError()),
		},
		{
			name: "err at ById",
			fields: fields{
				repo: NewMockRepo(&Config{
					ErrById: NewMockError(),
				}),
			},
			args: args{
				ctx:        context.TODO(),
				courseMeta: predefinedCourse.CourseMeta,
			},
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to retrieve the course: %w", NewMockError()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			courseManager, err := NewCourseManager(tt.fields.repo, tt.fields.logger)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			got, err := courseManager.Create(tt.args.ctx, tt.args.courseMeta)
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
	predefinedCourseWith10Students := generateUsersInCourse(10)
	predefinedCourseWith20Students := generateUsersInCourse(20)
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
				repo: NewMockRepo(&Config{
					ErrById: NewMockError(),
				}),
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
			name: "course max student validation",
			fields: fields{
				repo: NewMockRepo(&Config{CourseByUUID: map[uuid.UUID]models.Course{
					predefinedCourseWith20Students.Uuid: predefinedCourseWith20Students,
				}}),
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: predefinedCourseWith20Students.Uuid,
				student:    models.Student{},
			},
			wantErr:            true,
			expectedErrMessage: NewCourseConstraintErr(courseMaxStudentMsg).Error(),
		},
		{
			name: "err at ByStudent",
			fields: fields{
				repo: NewMockRepo(&Config{
					ErrByStudent: NewMockError(),
					CourseByUUID: map[uuid.UUID]models.Course{
						predefinedCourseWith10Students.Uuid: predefinedCourseWith10Students,
					},
				}),
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: predefinedCourseWith10Students.Uuid,
				student: models.Student{
					User: models.User{Uuid: fixedUuid},
				},
			},
			wantErr:            true,
			expectedErrMessage: fmt.Errorf("unable to retrieve courses: %w", NewMockError()).Error(),
		},
		{
			name: "error at Update",
			fields: fields{
				repo: NewMockRepo(&Config{
					CourseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourseWith10Students,
					},
					ErrUpdate: NewMockError(),
				}),
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
				repo: NewMockRepo(&Config{
					CourseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourseWith10Students,
					},
				}),
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
				t.Errorf("RegisterStudent() error = %v, wantErr %v", err.Error(), tt.expectedErrMessage)
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
					t.Errorf("failed to register the studentUUID. expected %v, got %v", tt.args.student.Uuid, courses[0].Students[tt.args.student.Uuid].Uuid)
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
		ctx         context.Context
		courseUUID  uuid.UUID
		studentUUID uuid.UUID
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
				repo: NewMockRepo(&Config{
					ErrById: NewMockError(),
				}),
			},
			args: args{
				ctx:         context.TODO(),
				courseUUID:  uuid.New(),
				studentUUID: uuid.New(),
			},
			wantErr:            true,
			expectedErrMessage: "unable to retrieve the course: mock error",
		},
		{
			name: "error at Update",
			fields: fields{
				repo: NewMockRepo(&Config{
					CourseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourse,
					},
					ErrUpdate: NewMockError(),
				}),
			},
			args: args{
				ctx:         context.TODO(),
				courseUUID:  fixedUuid,
				studentUUID: fixedUuid,
			},
			wantErr:            true,
			expectedErrMessage: "unable to update the course: mock error",
		},
		{
			name: "successful un-registry",
			fields: fields{
				repo: NewMockRepo(&Config{
					CourseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: predefinedCourse,
					},
				}),
			},
			args: args{
				ctx:         context.TODO(),
				courseUUID:  fixedUuid,
				studentUUID: anExistingStudent.Uuid,
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
			err = c.UnregisterStudent(tt.args.ctx, tt.args.courseUUID, tt.args.studentUUID)
			if tt.wantErr && tt.expectedErrMessage != err.Error() {
				t.Errorf("UnregisterStudent() error = %v, wantErr %v", err.Error(), tt.expectedErrMessage)
			}
			if !tt.wantErr {
				if err != nil {
					t.Errorf("unexpected error UnregisterStudent() error = %v", err)
				}
				courses, err := c.repo.ByStudent(tt.args.ctx, tt.args.studentUUID)
				if err != nil {
					t.Fatal("unexpected error", err)
				}
				for _, course := range courses {
					if course.Uuid == tt.args.courseUUID {
						t.Errorf("failed to un-register the studentUUID from the courseMeta %v", course)
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
				repo: NewMockRepo(&Config{
					ErrDelete: NewMockError(),
				}),
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: uuid.New(),
			},
			wantErr:            true,
			expectedErrMessage: "unable to delete the courseMeta: mock error",
		},
		{
			name: "successful Delete",
			fields: fields{
				repo: NewMockRepo(&Config{
					CourseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: {},
					},
				}),
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
				_, err := c.repo.ById(tt.args.ctx, tt.args.courseUUID)
				if err.Error() != "course not found" {
					t.Errorf("failed to delete course %v", tt.args.courseUUID)
				}

			}

		})
	}
}

func TestCourseManager_List(t *testing.T) {
	tutor := models.Tutor{
		User: models.User{
			Uuid: uuid.New(),
		},
	}
	type fields struct {
		repo   Repo
		logger *log.Logger
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name               string
		fields             fields
		args               args
		want               []models.Course
		wantErr            bool
		expectedErrMessage string
	}{
		{
			name: "error at repo list",
			fields: fields{
				repo: NewMockRepo(&Config{
					ErrList: NewMockError(),
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			wantErr:            true,
			expectedErrMessage: "unable to retrieve courses: mock error",
		},
		{
			name: "successful listing",
			fields: fields{
				repo: NewMockRepo(&Config{
					CourseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: {
							CourseMeta: models.CourseMeta{
								Uuid:  fixedUuid,
								Name:  "courseMeta 1",
								Tutor: &tutor,
							},
						},
					},
				}),
			},
			args: args{
				ctx: context.TODO(),
			},
			want: []models.Course{
				{CourseMeta: models.CourseMeta{
					Uuid:  fixedUuid,
					Name:  "courseMeta 1",
					Tutor: &tutor,
				},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCourseManager(tt.fields.repo, tt.fields.logger)
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			got, err := c.List(tt.args.ctx)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but none raised")
				}
				if tt.expectedErrMessage != err.Error() {
					t.Errorf("RegisterStudent() error = %v, wantErr %v", err.Error(), tt.expectedErrMessage)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error List() error = %v", err)
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("List() got = %v, want %v", got, tt.want)
				}
			}

		})
	}
}

func TestCourseManager_Get(t *testing.T) {
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
		want               *models.Course
		wantErr            bool
		expectedErrMessage string
	}{
		{
			name: "err at ById",
			fields: fields{
				repo: NewMockRepo(&Config{
					ErrById: NewMockError(),
				}),
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: uuid.New(),
			},
			wantErr:            true,
			expectedErrMessage: "unable to retrieve courseMeta by UUID: mock error",
		},
		{
			name: "successful Get",
			fields: fields{
				repo: NewMockRepo(&Config{
					CourseByUUID: map[uuid.UUID]models.Course{
						fixedUuid: {CourseMeta: models.CourseMeta{Uuid: fixedUuid}},
					},
				}),
			},
			args: args{
				ctx:        context.TODO(),
				courseUUID: fixedUuid,
			},
			want: &models.Course{CourseMeta: models.CourseMeta{Uuid: fixedUuid}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := NewCourseManager(tt.fields.repo, tt.fields.logger)
			if err != nil {
				t.Fatal("unexpected error", err)
			}
			got, err := c.Get(tt.args.ctx, tt.args.courseUUID)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, but none raised")
				}
				if tt.expectedErrMessage != err.Error() {
					t.Errorf("Get() error = %v, wantErr %v", err.Error(), tt.expectedErrMessage)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error Get() error = %v", err)
				}
				if !reflect.DeepEqual(*got, *tt.want) {
					t.Errorf("Get() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
