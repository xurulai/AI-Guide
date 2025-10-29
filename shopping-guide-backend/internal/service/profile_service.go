package service

import (
	"context"
	"shopping-guide-backend/internal/model"
	"shopping-guide-backend/internal/repository"
)

// ProfileService 用户画像服务接口
type ProfileService interface {
	GetProfile(ctx context.Context, userID string) (*model.UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, profile *model.UserProfile) error
}

// ProfileServiceImpl 用户画像服务实现
type ProfileServiceImpl struct {
	repo repository.ProfileRepository
}

func NewProfileService(repo repository.ProfileRepository) ProfileService {
	return &ProfileServiceImpl{
		repo: repo,
	}
}

func (s *ProfileServiceImpl) GetProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	return s.repo.GetProfile(ctx, userID)
}

func (s *ProfileServiceImpl) UpdateProfile(ctx context.Context, userID string, profile *model.UserProfile) error {
	return s.repo.UpdateProfile(ctx, userID, profile)
}
