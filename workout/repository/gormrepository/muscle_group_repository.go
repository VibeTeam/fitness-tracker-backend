package gormrepository

import (
	"context"

	"gorm.io/gorm"

	"github.com/VibeTeam/fitness-tracker-backend/workout/models"
	"github.com/VibeTeam/fitness-tracker-backend/workout/repository"
)

// gormMuscleGroupRepository implements repository.MuscleGroupRepository using GORM.
type gormMuscleGroupRepository struct {
	db *gorm.DB
}

// NewMuscleGroupRepository returns a GORM-backed MuscleGroup repository.
func NewMuscleGroupRepository(db *gorm.DB) repository.MuscleGroupRepository {
	return &gormMuscleGroupRepository{db: db}
}

func (r *gormMuscleGroupRepository) Create(ctx context.Context, mg *models.MuscleGroup) error {
	return r.db.WithContext(ctx).Create(mg).Error
}

func (r *gormMuscleGroupRepository) GetByID(ctx context.Context, id uint) (*models.MuscleGroup, error) {
	var mg models.MuscleGroup
	err := r.db.WithContext(ctx).First(&mg, id).Error
	if err != nil {
		return nil, err
	}
	return &mg, nil
}

func (r *gormMuscleGroupRepository) Update(ctx context.Context, mg *models.MuscleGroup) error {
	return r.db.WithContext(ctx).Save(mg).Error
}

func (r *gormMuscleGroupRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.MuscleGroup{}, id).Error
}

func (r *gormMuscleGroupRepository) List(ctx context.Context, limit, offset int) ([]*models.MuscleGroup, error) {
	var mgs []*models.MuscleGroup
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&mgs).Error
	return mgs, err
}

func (r *gormMuscleGroupRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.MuscleGroup{}).Count(&count).Error
	return int(count), err
}
