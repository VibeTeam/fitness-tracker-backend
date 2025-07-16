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
