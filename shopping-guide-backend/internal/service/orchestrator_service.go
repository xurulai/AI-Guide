package service

import (
	"context"
	"fmt"

	"shopping-guide-backend/internal/model"
)

// OrchestratorService 编排服务
// 职责：协调Planner和Executor的调用流程
type OrchestratorService interface {
	// ProcessChat 处理对话（核心流程）
	// 1. 加载会话上下文
	// 2. 调用Planner进行规划
	// 3. 调用对应的Executor执行
	// 4. 保存会话和日志
	ProcessChat(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error)
}

// orchestratorService 编排服务实现
type orchestratorService struct {
	plannerService  PlannerService
	executorService ExecutorService
	sessionService  SessionService
	productService  ProductService
	profileService  ProfileService
	// logRepo repository.LogRepository
}

// NewOrchestratorService 创建编排服务
func NewOrchestratorService(
	plannerService PlannerService,
	executorService ExecutorService,
	sessionService SessionService,
	productService ProductService,
	profileService ProfileService,
) OrchestratorService {
	return &orchestratorService{
		plannerService:  plannerService,
		executorService: executorService,
		sessionService:  sessionService,
		productService:  productService,
		profileService:  profileService,
	}
}

// ProcessChat 处理对话
func (s *orchestratorService) ProcessChat(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error) {
	// TODO: 实现编排逻辑
	// 步骤：
	// 1. 获取会话上下文: session := s.sessionService.GetSession(...)
	// 2. 获取商品库: productStorage := s.productService.GetProductStorage(...)
	// 3. 调用Planner: plannerResult := s.plannerService.Analyze(...)
	// 4. 调用Executor: executorResult := s.executorService.Execute(...)
	// 5. 保存会话: s.sessionService.SaveSession(...)
	// 6. 返回响应
	session, err := s.sessionService.GetSession(ctx, req.SessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		session, err = s.sessionService.CreateSession(ctx, &model.SessionCreateRequest{
			UserID: req.UserID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create session: %w", err)
		}
	}

	plannerResult, err := s.plannerService.Analyze(ctx, &PlannerRequest{
		Query:     req.Query,
		History:   session.Messages,
		SessionID: session.SessionID,
		UserID:    session.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze: %w", err)
	}

	session.UserID = "user_001"

	userProfile, err := s.profileService.GetProfile(ctx, session.UserID)
	if err != nil {
		// 如果获取失败，使用默认画像
		fmt.Printf("⚠️  Failed to get user profile: %v, using default profile\n", err)
		userProfile = &model.UserProfile{
			UserID:         session.UserID,
			PreferredStyle: "xiaohongshu",
			Age:            25,
			Gender:         "unknown",
			Interests:      "[]", // 空的 JSON 数组
		}
	}

	// 根据Planner结果选择对应的Executor
	executorReq := &ExecutorRequest{
		Query:               req.Query,
		Tool:                plannerResult.Tool,
		UserProfile:         *userProfile,
		History:             session.Messages,
		BusinessInstruction: session.BusinessInstruction,
		UserID:              session.UserID,
	}

	executorResult, err := s.executorService.Execute(ctx, executorReq)
	if err != nil {
		return nil, fmt.Errorf("executor execute failed: %w", err)
	}

	return &model.ChatResponse{
		SessionID: session.SessionID,
		Response:  executorResult.Response,
		ToolUsed:  plannerResult.Tool,
	}, nil

}
