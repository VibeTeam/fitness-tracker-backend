package gormrepository

import (
	"context"

	"gorm.io/gorm"

	"github.com/VibeTeam/fitness-tracker-backend/workout/models"
	"github.com/VibeTeam/fitness-tracker-backend/workout/repository"
)

// gormWorkoutDetailRepository implements repository.WorkoutDetailRepository using GORM.
type gormWorkoutDetailRepository struct {
	db *gorm.DB
}

// NewWorkoutDetailRepository returns a GORM-backed WorkoutDetail repository.
func NewWorkoutDetailRepository(db *gorm.DB) repository.WorkoutDetailRepository {
	return &gormWorkoutDetailRepository{db: db}
}

func (r *gormWorkoutDetailRepository) Create(ctx context.Context, detail *models.WorkoutDetail) error {
	return r.db.WithContext(ctx).Create(detail).Error
}

func (r *gormWorkoutDetailRepository) GetByID(ctx context.Context, id uint) (*models.WorkoutDetail, error) {
	var detail models.WorkoutDetail
	err := r.db.WithContext(ctx).First(&detail, id).Error
	if err != nil {
		return nil, err
	}
	return &detail, nil
}

func (r *gormWorkoutDetailRepository) Update(ctx context.Context, detail *models.WorkoutDetail) error {
	return r.db.WithContext(ctx).Save(detail).Error
}

func (r *gormWorkoutDetailRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.WorkoutDetail{}, id).Error
}

func (r *gormWorkoutDetailRepository) ListBySession(ctx context.Context, sessionID uint) ([]*models.WorkoutDetail, error) {
	var details []*models.WorkoutDetail
	err := r.db.WithContext(ctx).Where("workout_session_id = ?", sessionID).Find(&details).Error
	return details, err
}
