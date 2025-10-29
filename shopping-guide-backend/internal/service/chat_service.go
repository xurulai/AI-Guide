package service

import (
	"context"

	"shopping-guide-backend/internal/model"

	"github.com/google/uuid"
)

// ChatService 对话服务接口（对外暴露的顶层服务）
type ChatService interface {
	Chat(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error)
	ChatStream(ctx context.Context, req *model.ChatRequest) (<-chan model.StreamChunk, error)
}

// chatService 对话服务实现
type chatService struct {
	orchestrator OrchestratorService
}

// NewChatService 创建对话服务
func NewChatService(orchestrator OrchestratorService) ChatService {
	return &chatService{
		orchestrator: orchestrator,
	}
}

// Chat 对话（阻塞式）
func (s *chatService) Chat(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error) {

	sessionID := uuid.New().String()
	req.SessionID = sessionID
	return s.orchestrator.ProcessChat(ctx, req)
}

// ChatStream 对话（流式）
func (s *chatService) ChatStream(ctx context.Context, req *model.ChatRequest) (<-chan model.StreamChunk, error) {
	// TODO: 实现流式对话
	ch := make(chan model.StreamChunk)
	close(ch)
	return ch, nil
}
