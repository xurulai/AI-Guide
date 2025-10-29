package model

import "time"

// Session 会话模型
type Session struct {
	SessionID           string                 `json:"session_id"`
	UserID              string                 `json:"user_id"`
	Messages            []Message              `json:"messages"`
	ProductStorage      map[string]interface{} `json:"product_storage"`
	BusinessInstruction string                 `json:"business_instruction"`
	UserProfile         UserProfile            `json:"user_profile"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	ExpiresAt           time.Time              `json:"expires_at"`
}

// Message 消息模型
type Message struct {
	Role      string                 `json:"role"` // user/assistant
	Content   string                 `json:"content"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AddMessage 添加消息到会话
func (s *Session) AddMessage(role, content string) {
	s.Messages = append(s.Messages, Message{
		Role:      role,
		Content:   content,
		Timestamp: time.Now(),
	})
	s.UpdatedAt = time.Now()
}

// GetRecentMessages 获取最近N条消息
func (s *Session) GetRecentMessages(n int) []Message {
	if len(s.Messages) <= n {
		return s.Messages
	}
	return s.Messages[len(s.Messages)-n:]
}

// IsExpired 检查会话是否过期
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
