package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"shopping-guide-backend/internal/config"
	"shopping-guide-backend/internal/model"

	"github.com/go-redis/redis/v8"
)

// SessionRepository 会话存储接口
type SessionRepository interface {
	Save(ctx context.Context, session *model.Session) error
	Get(ctx context.Context, sessionID string) (*model.Session, error)
	Delete(ctx context.Context, sessionID string) error
	AppendMessage(ctx context.Context, sessionID string, message *model.Message) error
	GetMessages(ctx context.Context, sessionID string) ([]model.Message, error)
	GetRecentMessages(ctx context.Context, sessionID string, n int) ([]model.Message, error)
	ClearMessages(ctx context.Context, sessionID string) error
}

type sessionRepository struct {
	rdb *redis.Client
	cfg *config.RedisConfig
}

func NewSessionRepository(rdb *redis.Client, cfg *config.RedisConfig) SessionRepository {
	return &sessionRepository{
		rdb: rdb,
		cfg: cfg,
	}
}
func (s sessionRepository) Save(ctx context.Context, session *model.Session) error {

	key := fmt.Sprintf("session:%s", session.SessionID)

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	err = s.rdb.Set(ctx, key, data, s.cfg.SessionTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to set session: %w", err)
	}
	return nil
}

// Get 获取会话数据
func (s sessionRepository) Get(ctx context.Context, sessionID string) (*model.Session, error) {

	key := fmt.Sprintf("session:%s", sessionID)

	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	var session model.Session

	if err := json.Unmarshal([]byte(data), &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// Delete 删除会话数据 (包括消息历史)
func (s sessionRepository) Delete(ctx context.Context, sessionID string) error {
	sessionkey := fmt.Sprintf("session:%s", sessionID)
	messagekey := fmt.Sprintf("session:%s:messages", sessionID)
	err := s.rdb.Del(ctx, sessionkey, messagekey).Err()
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}
	return nil
}

// AppendMessage 追加单条消息
func (s sessionRepository) AppendMessage(ctx context.Context, sessionID string, message *model.Message) error {
	messagekey := fmt.Sprintf("session:%s:messages", sessionID)

	//序列化消息
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	//追加消息
	err = s.rdb.RPush(ctx, messagekey, data).Err()
	if err != nil {
		return fmt.Errorf("failed to append message: %w", err)
	}

	err = s.rdb.Expire(ctx, messagekey, s.cfg.SessionTTL).Err()
	if err != nil {
		return fmt.Errorf("failed to expire message: %w", err)
	}
	return nil
}

// GetMessages 获取所有消息历史
func (s sessionRepository) GetMessages(ctx context.Context, sessionID string) ([]model.Message, error) {
	messagekey := fmt.Sprintf("session:%s:messages", sessionID)

	vals, err := s.rdb.LRange(ctx, messagekey, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get messages: %w", err)
	}

	messages := make([]model.Message, len(vals))
	for _, val := range vals {
		var message model.Message
		if err := json.Unmarshal([]byte(val), &message); err != nil {
			return nil, fmt.Errorf("failed to unmarshal message: %w", err)
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// GetRecentMessages 获取最近 N 条消息
func (s sessionRepository) GetRecentMessages(ctx context.Context, sessionID string, n int) ([]model.Message, error) {
	messagekey := fmt.Sprintf("session:%s:messages", sessionID)

	// 获取列表长度
	length, err := s.rdb.LLen(ctx, messagekey).Result()
	if err != nil {
		if err == redis.Nil {
			return []model.Message{}, nil
		}
		return nil, fmt.Errorf("llen failed: %w", err)
	}

	if length == 0 {
		return []model.Message{}, nil
	}

	// 计算起始位置（取最后 n 条）
	start := int64(0)
	if int64(n) < length {
		start = length - int64(n)
	}

	// 获取指定范围的消息
	vals, err := s.rdb.LRange(ctx, messagekey, start, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("lrange failed: %w", err)
	}

	// 反序列化
	messages := make([]model.Message, 0, len(vals))
	for _, s := range vals {
		var m model.Message
		if e := json.Unmarshal([]byte(s), &m); e == nil {
			messages = append(messages, m)
		}
	}

	return messages, nil
}

// ClearMessages 清空会话的所有消息历史
func (s sessionRepository) ClearMessages(ctx context.Context, sessionID string) error {
	messagekey := fmt.Sprintf("session:%s:messages", sessionID)
	err := s.rdb.Del(ctx, messagekey).Err()
	if err != nil {
		return fmt.Errorf("failed to clear messages: %w", err)
	}
	return nil
}
