package repository

import (
	"context"

	"github.com/VibeTeam/fitness-tracker-backend/workout/models"
)

// MuscleGroupRepository provides CRUD operations for MuscleGroup entities.
type MuscleGroupRepository interface {
	Create(ctx context.Context, mg *models.MuscleGroup) error
	GetByID(ctx context.Context, id uint) (*models.MuscleGroup, error)
	Update(ctx context.Context, mg *models.MuscleGroup) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*models.MuscleGroup, error)
	Count(ctx context.Context) (int, error)
}

// WorkoutTypeRepository provides CRUD operations for WorkoutType entities.
type WorkoutTypeRepository interface {
	Create(ctx context.Context, wt *models.WorkoutType) error
	GetByID(ctx context.Context, id uint) (*models.WorkoutType, error)
	Update(ctx context.Context, wt *models.WorkoutType) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, limit, offset int) ([]*models.WorkoutType, error)
	Count(ctx context.Context) (int, error)
}

// WorkoutSessionRepository provides operations for WorkoutSession and its details.
type WorkoutSessionRepository interface {
	Create(ctx context.Context, session *models.WorkoutSession) error
	GetByID(ctx context.Context, id uint) (*models.WorkoutSession, error)
	Update(ctx context.Context, session *models.WorkoutSession) error
	Delete(ctx context.Context, id uint) error

	// List all sessions for a specific user with pagination.
	ListByUser(ctx context.Context, userID uint, limit, offset int) ([]*models.WorkoutSession, error)
	CountByUser(ctx context.Context, userID uint) (int, error)
}

// WorkoutDetailRepository provides CRUD operations for WorkoutDetail entities.
type WorkoutDetailRepository interface {
	Create(ctx context.Context, detail *models.WorkoutDetail) error
	GetByID(ctx context.Context, id uint) (*models.WorkoutDetail, error)
	Update(ctx context.Context, detail *models.WorkoutDetail) error
	Delete(ctx context.Context, id uint) error
	ListBySession(ctx context.Context, sessionID uint) ([]*models.WorkoutDetail, error)
}
