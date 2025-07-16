package repository

import (
	"context"

	"github.com/VibeTeam/fitness-tracker-backend/workout/models"
)

// WorkoutDetailRepository provides CRUD operations for WorkoutDetail entities.
type WorkoutDetailRepository interface {
	Create(ctx context.Context, detail *models.WorkoutDetail) error
	GetByID(ctx context.Context, id uint) (*models.WorkoutDetail, error)
	Update(ctx context.Context, detail *models.WorkoutDetail) error
	Delete(ctx context.Context, id uint) error
	ListBySession(ctx context.Context, sessionID uint) ([]*models.WorkoutDetail, error)
}