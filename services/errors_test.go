package services

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestCourseConstraintErr_Error(t *testing.T) {
	type fields struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "valid error",
			fields: fields{message: "violated some constraint"},
			want:   "violated some constraint",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &CourseConstraintErr{
				message: tt.fields.message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCourseConstraintErr_Is(t *testing.T) {
	type fields struct {
		message string
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "expected error",
			fields: fields{message: fmt.Sprintf(validationErrFmt, string(tutorMaxCourseMsg))},
			args:   args{target: NewCourseConstraintErr(tutorMaxCourseMsg)},
			want:   true,
		},
		{
			name:   "not expected error",
			fields: fields{message: fmt.Sprintf(validationErrFmt, "random string")},
			args:   args{target: NewCourseConstraintErr(tutorMaxCourseMsg)},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &CourseConstraintErr{
				message: tt.fields.message,
			}
			if got := e.Is(tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCourseConstraintErr(t *testing.T) {
	type args struct {
		check courseConstraint
	}
	tests := []struct {
		name string
		args args
		want *CourseConstraintErr
	}{
		{
			name: "valid error tutor max courseMeta",
			args: args{check: tutorMaxCourseMsg},
			want: &CourseConstraintErr{message: fmt.Sprintf(validationErrFmt, tutorMaxCourseMsg)},
		},
		{
			name: "valid error studentUUID max courseMeta",
			args: args{check: studentMaxCourseMsg},
			want: &CourseConstraintErr{message: fmt.Sprintf(validationErrFmt, studentMaxCourseMsg)},
		},
		{
			name: "valid error courseMeta max studentUUID",
			args: args{check: courseMaxStudentMsg},
			want: &CourseConstraintErr{message: fmt.Sprintf(validationErrFmt, courseMaxStudentMsg)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCourseConstraintErr(tt.args.check); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCourseConstraintErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNilErr(t *testing.T) {
	type args struct {
		item string
	}
	tests := []struct {
		name string
		args args
		want *NilErr
	}{
		{
			name: "valid error",
			args: args{item: "some item"},
			want: NewNilErr("some item"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNilErr(tt.args.item); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNilErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewNotFoundErr(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name string
		args args
		want *NotFoundError
	}{
		{
			name: "valid error",
			args: args{message: "some item"},
			want: NewNotFoundErr("some item"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotFoundErr(tt.args.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotFoundErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCourseNotFoundErr(t *testing.T) {
	type args struct {
		courseUuid uuid.UUID
	}
	tests := []struct {
		name string
		args args
		want *NotFoundError
	}{
		{
			name: "valid error",
			args: args{courseUuid: uuid.MustParse("5d61cbc8-9ccd-4348-a623-d61dd7658dd7")},
			want: NewCourseNotFoundErr(uuid.MustParse("5d61cbc8-9ccd-4348-a623-d61dd7658dd7")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCourseNotFoundErr(tt.args.courseUuid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCourseNotFoundErr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilErr_Error(t *testing.T) {
	type fields struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "valid error",
			fields: fields{message: "this cannot be nil"},
			want:   "this cannot be nil",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &NilErr{
				message: tt.fields.message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNilErr_Is(t *testing.T) {
	type fields struct {
		message string
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "matching error",
			fields: fields{message: "this cannot be nil"},
			args:   args{target: NewNilErr("this")},
			want:   true,
		},
		{
			name:   "not-matching error",
			fields: fields{message: "this cannot be nil"},
			args:   args{target: NewNilErr("that")},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &NilErr{
				message: tt.fields.message,
			}
			if got := e.Is(tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFoundError_Error(t *testing.T) {
	type fields struct {
		message string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "valid error",
			fields: fields{message: "this is not found"},
			want:   "this is not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &NotFoundError{
				message: tt.fields.message,
			}
			if got := e.Error(); got != tt.want {
				t.Errorf("Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotFoundError_Is(t *testing.T) {
	type fields struct {
		message string
	}
	type args struct {
		target error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name:   "matching error",
			fields: fields{message: "this is not found"},
			args:   args{target: NewNotFoundErr("this is not found")},
			want:   true,
		},
		{
			name:   "un-matching error",
			fields: fields{message: "this is not found"},
			args:   args{target: NewNotFoundErr("that is not found")},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &NotFoundError{
				message: tt.fields.message,
			}
			if got := e.Is(tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
