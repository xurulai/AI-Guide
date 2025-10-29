package service

import (
	"context"
	"fmt"
	"shopping-guide-backend/internal/config"
	"shopping-guide-backend/internal/model"
	"shopping-guide-backend/internal/repository"
	"time"

	"github.com/google/uuid"
)

// SessionService 会话服务接口
type SessionService interface {
	CreateSession(ctx context.Context, req *model.SessionCreateRequest) (*model.Session, error)
	GetSession(ctx context.Context, sessionID string) (*model.Session, error)
	SaveSession(ctx context.Context, session *model.Session) error
	//DeleteSession(ctx context.Context, sessionID string) error
}

type sessionServiceImpl struct {
	repo        repository.SessionRepository
	productRepo repository.ProductRepository
	cfg         *config.Config
}

// NewSessionServiceImpl 创建会话服务实现
func NewSessionServiceImpl(
	repo repository.SessionRepository,
	productRepo repository.ProductRepository,
	cfg *config.Config,
) SessionService {
	return &sessionServiceImpl{
		repo:        repo,
		productRepo: productRepo,
		cfg:         cfg,
	}
}

// CreateSession 创建新会话
func (s *sessionServiceImpl) CreateSession(ctx context.Context, req *model.SessionCreateRequest) (*model.Session, error) {

	sessionID := uuid.New().String()

	now := time.Now()

	session := &model.Session{
		SessionID:           sessionID,
		UserID:              req.UserID,
		Messages:            []model.Message{}, // 初始为空
		BusinessInstruction: req.BusinessInstruction,
		UserProfile: model.UserProfile{
			PreferredStyle: s.cfg.Business.Session.DefaultStyle,
		},
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: now.Add(s.cfg.Redis.SessionTTL),
	}

	//保存会话数据
	if err := s.repo.Save(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to save session: %w", err)
	}

	return session, nil

}

func (s *sessionServiceImpl) GetSession(ctx context.Context, sessionID string) (*model.Session, error) {
	session, err := s.repo.Get(ctx, sessionID)
	if err != nil {
		if err.Error() == "session not found" {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

func (s *sessionServiceImpl) SaveSession(ctx context.Context, session *model.Session) error {
	return s.repo.Save(ctx, session)
}
