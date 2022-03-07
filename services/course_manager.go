package services

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/tomasdembelli/course-manager/models"
	"log"
)

// Repo is the interface that defines the methods for persisting and manipulating service data.
type Repo interface {
	ById(ctx context.Context, courseUuid uuid.UUID) (models.Course, error)
	ByTutor(ctx context.Context, tutorUuid uuid.UUID) ([]models.Course, error)
	ByStudent(ctx context.Context, studentUuid uuid.UUID) ([]models.Course, error)
	List(ctx context.Context) ([]models.Course, error)
	Create(ctx context.Context, course models.Course) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	Update(ctx context.Context, uuid uuid.UUID, course models.Course) error
}

// Course is the service for managing the courses.
type Course struct {
	repo   Repo
	logger *log.Logger
}

// NewCourse initiates a new course service with the given repo.
func NewCourse(repo Repo, logger *log.Logger) (Course, error) {
	if repo == nil {
		return Course{}, NewNilErr("repo")
	}

	if logger == nil {
		return Course{
			repo:   repo,
			logger: log.Default(),
		}, nil
	}

	return Course{
		repo:   repo,
		logger: logger,
	}, nil
}

// Create creates a new course after validating the following checks:
//	- A tutor can facilitate maximum 2 courses.
//	- A student can register to maximum 4 courses.
//	- Maximum 20 students can register a course.
func (c *Course) Create(ctx context.Context, course models.Course) (*models.Course, error) {
	coursesByTutor, err := c.repo.ByTutor(ctx, course.Tutor.Uuid)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve courses: %w", err)
	}
	if len(coursesByTutor) >= 2 {
		return nil, NewCourseConstraintErr(tutorMaxCourse)
	}
	if len(course.Students) > 0 {
		if len(course.Students) > 20 {
			return nil, NewCourseConstraintErr(courseMaxStudent)
		}
		var notEligibleStudents []uuid.UUID
		for _, student := range course.Students {
			coursesByStudent, err := c.repo.ByStudent(ctx, student.Uuid)
			if err != nil {
				return nil, fmt.Errorf("unable to retrieve courses: %w", err)
			}
			if len(coursesByStudent) >= 4 {
				notEligibleStudents = append(notEligibleStudents, student.Uuid)
				delete(course.Students, student.Uuid)
			}
		}
		if len(notEligibleStudents) > 0 {
			c.logger.Printf("%v, ineligible student uuids: %v", studentMaxCourse, notEligibleStudents)
		}
	}
	courseUuid := uuid.New()
	course.Uuid = courseUuid
	err = c.repo.Create(ctx, course)
	if err != nil {
		return nil, fmt.Errorf("unable to create the course: %w", err)
	}

	courseCreated, err := c.repo.ById(ctx, courseUuid)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the course: %w", err)
	}

	return &courseCreated, nil
}
