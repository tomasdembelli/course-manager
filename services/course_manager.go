package services

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/tomasdembelli/course-manager/models"
)

// Repo is the interface that defines the methods for persisting and manipulating service data.
type Repo interface {
	ById(ctx context.Context, courseUUID uuid.UUID) (*models.Course, error)
	ByTutor(ctx context.Context, tutorUUID uuid.UUID) ([]models.Course, error)
	ByStudent(ctx context.Context, studentUUID uuid.UUID) ([]models.Course, error)
	List(ctx context.Context) ([]models.Course, error)
	Create(ctx context.Context, course models.Course) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	Update(ctx context.Context, course models.Course) error
}

// CourseManager is the service for managing the courses.
type CourseManager struct {
	repo   Repo
	logger *log.Logger
}

// NewCourseManager initiates a new CourseManager service with the given repo.
func NewCourseManager(repo Repo, logger *log.Logger) (CourseManager, error) {
	if repo == nil {
		return CourseManager{}, NewNilErr("repo")
	}

	if logger == nil {
		return CourseManager{
			repo:   repo,
			logger: log.Default(),
		}, nil
	}

	return CourseManager{
		repo:   repo,
		logger: logger,
	}, nil
}

// Create creates a new course after validating the following checks:
//	- A tutor can facilitate maximum 2 courses.
//	- A student can register to maximum 4 courses.
//	- Maximum 20 students can register a course.
func (c *CourseManager) Create(ctx context.Context, course models.Course) (*models.Course, error) {
	if course.Tutor == nil {
		return nil, NewNilErr("tutor")
	}
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
	if course.Uuid == uuid.Nil {
		course.Uuid = uuid.New()
	}

	err = c.repo.Create(ctx, course)
	if err != nil {
		return nil, fmt.Errorf("unable to create the course: %w", err)
	}

	courseCreated, err := c.repo.ById(ctx, course.Uuid)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve the course: %w", err)
	}

	return courseCreated, nil
}

// RegisterStudent registers the given models.Student to the given course.
// This is an idempotent operation.
// It will return an error if the given course is not found or unable to update it.
func (c CourseManager) RegisterStudent(ctx context.Context, courseUUID uuid.UUID, student models.Student) error {
	course, err := c.repo.ById(ctx, courseUUID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the course: %w", err)
	}
	course.Students[student.Uuid] = student
	err = c.repo.Update(ctx, *course)
	if err != nil {
		return fmt.Errorf("unable to update the course: %w", err)
	}
	return nil
}

// UnregisterStudent removes the given models.Student from the given course.
// This is an idempotent operation.
// It will return an error if the given course is not found or unable to update it.
// If the student has not been registered to the course previously, no error will be returned (no-op).
func (c CourseManager) UnregisterStudent(ctx context.Context, courseUUID, studentUUID uuid.UUID) error {
	course, err := c.repo.ById(ctx, courseUUID)
	if err != nil {
		return fmt.Errorf("unable to retrieve the course: %w", err)
	}
	delete(course.Students, studentUUID)
	err = c.repo.Update(ctx, *course)
	if err != nil {
		return fmt.Errorf("unable to update the course: %w", err)
	}
	return nil
}

// Delete deletes the course for the given courseUUID.
// This is an idempotent operation.
func (c *CourseManager) Delete(ctx context.Context, courseUUID uuid.UUID) error {
	err := c.repo.Delete(ctx, courseUUID)
	if err != nil {
		return fmt.Errorf("unable to delete the course: %w", err)
	}
	return nil
}

// List returns all courses in the repo.
// It returns an error if it fails to fetch courses from the repo.
func (c *CourseManager) List(ctx context.Context) ([]models.Course, error) {
	courses, err := c.repo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve courses: %w", err)
	}
	return courses, nil
}

// Get returns the models.Course for the given course UUID.
func (c CourseManager) Get(ctx context.Context, courseUUID uuid.UUID) (*models.Course, error) {
	course, err := c.repo.ById(ctx, courseUUID)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve course by UUID: %w", err)
	}
	fmt.Println(course)
	if course.Uuid == uuid.Nil {
		return nil, NewCourseNotFoundErr(courseUUID)
	}
	return course, nil
}
