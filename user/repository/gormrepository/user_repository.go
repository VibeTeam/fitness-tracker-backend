package gormrepository

import (
	"context"

	"gorm.io/gorm"

	"github.com/VibeTeam/fitness-tracker-backend/user/models"
	"github.com/VibeTeam/fitness-tracker-backend/user/repository"
)

// gormUserRepository implements repository.UserRepository using GORM.
type gormUserRepository struct {
	db *gorm.DB
}

// NewUserRepository returns a GORM-backed User repository.
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &gormUserRepository{db: db}
}

func (r *gormUserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *gormUserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *gormUserRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *gormUserRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.User{}, id).Error
}

func (r *gormUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	var users []*models.User
	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error
	return users, err
}

func (r *gormUserRepository) Count(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.User{}).Count(&count).Error
	return int(count), err
}
