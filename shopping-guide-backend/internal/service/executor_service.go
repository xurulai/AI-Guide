package service

import (
	"context"
	"fmt"

	"shopping-guide-backend/internal/model"
)

// ExecutorService Executor执行器服务（从Agent）
// 根据Planner的Tool选择，调用对应的Executor
type ExecutorService interface {
	// Execute 执行具体的Agent逻辑
	Execute(ctx context.Context, req *ExecutorRequest) (*model.ExecutorResult, error)
}

// ExecutorRequest Executor请求
type ExecutorRequest struct {
	Query               string            // 用户输入
	Tool                string            // Tool类型（PRODUCT_RECOMMENDATION_MODULE等）
	History             []model.Message   // 对话历史
	BusinessInstruction string            // 商家场景描述
	UserProfile         model.UserProfile // 用户画像
	UserID              string            // 用户ID
}

// executorService Executor服务实现
type executorService struct {
	// TODO: 注入DifyClient
	// TODO: 注入ProductService
	// TODO: 注入LogRepository
}

// NewExecutorService 创建Executor服务
func NewExecutorService() ExecutorService {
	return &executorService{}
}

// Execute 执行
func (s *executorService) Execute(ctx context.Context, req *ExecutorRequest) (*model.ExecutorResult, error) {
	// TODO: 根据Tool类型选择对应的Executor工作流
	// TODO: 调用Dify Executor工作流
	// TODO: 解析返回结果
	// TODO: 记录日志
	switch req.Tool {
	case model.ToolShoppingGuide:
		return s.executeShoppingGuide(ctx, req)
	case model.ToolProductRecommendation:
		return s.executeProductRecommendation(ctx, req)
	case model.ToolQAAssistant:
		return s.executeQAAssistant(ctx, req)
	}

	return nil, nil
}

func (s *executorService) executeShoppingGuide(ctx context.Context, req *ExecutorRequest) (*model.ExecutorResult, error) {

	// 将 UserProfile 序列化为 JSON 字符串
	userPortraitJSON := fmt.Sprintf(`{"age": %d, "gender": "%s", "interests": %s}`,
		req.UserProfile.Age,
		req.UserProfile.Gender,
		req.UserProfile.Interests, // 已经是 JSON 字符串格式
	)

	// 构造请求参数
	inputs := map[string]interface{}{
		"query":         req.Query,
		"user_portrait": userPortraitJSON,
	}

	// 只有当 history 非空时才传递
	if len(req.History) > 0 {
		inputs["history"] = req.History
	}

	difyreq := model.DifyWorkflowRequest{
		Inputs:       inputs,
		ResponseMode: "blocking",
		User:         req.UserID,
	}

	fmt.Printf("\n🔍 Calling GuideChat with query: %s\n", req.Query)
	fmt.Printf("📊 User Portrait: %s\n", userPortraitJSON)

	difyresp, err := model.GuideChat(ctx, difyreq)
	if err != nil {
		return nil, fmt.Errorf("failed to call guide chat: %w", err)
	}

	fmt.Printf("\n📝 Guide raw response: %s\n", difyresp)

	// TODO: 解析 difyresp 并构造 ExecutorResult
	result := &model.ExecutorResult{
		Response:            difyresp,
		RecommendedProducts: []model.RecommendedProduct{},
		Metadata:            map[string]interface{}{},
	}

	return result, nil
}

func (s *executorService) executeProductRecommendation(ctx context.Context, req *ExecutorRequest) (*model.ExecutorResult, error) {
	return nil, nil
}

func (s *executorService) executeQAAssistant(ctx context.Context, req *ExecutorRequest) (*model.ExecutorResult, error) {
	return nil, nil
}
