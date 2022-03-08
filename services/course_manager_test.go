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

const mockErrMessage = "mock error"

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
	content map[uuid.UUID]models.Course
	err     error
}

func (m *mockRepo) ById(ctx context.Context, courseUuid uuid.UUID) (models.Course, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) ByTutor(_ context.Context, tutorUuid uuid.UUID) ([]models.Course, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, nil
}

func (m *mockRepo) ByStudent(ctx context.Context, studentUuid uuid.UUID) ([]models.Course, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) List(ctx context.Context) ([]models.Course, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) Create(ctx context.Context, course models.Course) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) Delete(ctx context.Context, uuid uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (m *mockRepo) Update(ctx context.Context, uuid uuid.UUID, course models.Course) error {
	//TODO implement me
	panic("implement me")
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

func TestCourseManager_Create(t *testing.T) {
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
			name: "error at ByTutor",
			fields: fields{repo: &mockRepo{
				content: nil,
				err:     newMockError(),
			}},
			args:        args{ctx: context.TODO(), course: models.Course{}},
			want:        nil,
			wantErr:     true,
			expectedErr: fmt.Errorf("unable to retrieve courses: %w", newMockError()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CourseManager{
				repo:   tt.fields.repo,
				logger: tt.fields.logger,
			}
			got, err := c.Create(tt.args.ctx, tt.args.course)
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
