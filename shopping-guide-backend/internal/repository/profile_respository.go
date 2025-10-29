package repository

import (
	"context"
	"shopping-guide-backend/internal/config"
	"shopping-guide-backend/internal/model"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context, userID string) (*model.UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, profile *model.UserProfile) error
}

type profileRepository struct {
	db  *gorm.DB
	cfg *config.MySQLConfig
}

func NewProfileRepository(db *gorm.DB, cfg *config.MySQLConfig) ProfileRepository {
	return &profileRepository{
		db:  db,
		cfg: cfg,
	}
}

func (p *profileRepository) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	var profile model.UserProfile
	if err := p.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}
	return &profile, nil
}

func (p *profileRepository) UpdateProfile(ctx context.Context, userID string, profile *model.UserProfile) error {
	return p.db.WithContext(ctx).Where("user_id = ?", userID).Save(profile).Error
}
