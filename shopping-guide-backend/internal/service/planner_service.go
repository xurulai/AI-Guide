package service

import (
	"context"
	"encoding/json"
	"fmt"

	"shopping-guide-backend/internal/model"
)

// PlannerService Planner规划器服务（Master Agent）
// 职责：意图识别、Tool选择、执行策略规划
type PlannerService interface {
	// Analyze 分析用户输入，返回规划结果
	// 包括：是否有明确购物意图、应该调用哪个Tool
	Analyze(ctx context.Context, req *PlannerRequest) (*model.PlannerResult, error)
}

// PlannerRequest Planner请求
type PlannerRequest struct {
	Query     string          // 用户输入
	History   []model.Message // 对话历史
	SessionID string          // 会话ID
	UserID    string          // 用户ID
}

// plannerService Planner服务实现
type plannerService struct {
	// TODO: 注入DifyClient
	// TODO: 注入LogRepository
}

// NewPlannerService 创建Planner服务
func NewPlannerService() PlannerService {
	return &plannerService{}
}

// Analyze 分析并规划
func (s *plannerService) Analyze(ctx context.Context, req *PlannerRequest) (*model.PlannerResult, error) {
	// 调用Dify Planner工作流
	difyreq := model.DifyWorkflowRequest{
		Inputs: map[string]interface{}{
			"query": req.Query,
		},
		ResponseMode: "blocking", // 使用 blocking 模式便于调试
		User:         req.UserID,
	}

	difyresp, err := model.PlannerChat(ctx, difyreq)
	if err != nil {
		return nil, fmt.Errorf("failed to call planner chat: %w", err)
	}

	fmt.Printf("\n📝 Planner raw response: %s\n", difyresp)

	// Dify 返回的可能是 JSON 数组，如 ["SHOPPING_GUIDE_AND_INTENT_MINING_MODULE"]
	// 需要先解析成数组，再提取第一个元素
	var tools []string
	if err := json.Unmarshal([]byte(difyresp), &tools); err != nil {
		// 如果不是数组，可能是直接返回的字符串，尝试作为单个 tool 处理
		return &model.PlannerResult{
			Tool: difyresp,
		}, nil
	}

	if len(tools) == 0 {
		return nil, fmt.Errorf("planner returned empty tools array")
	}

	plannerResult := &model.PlannerResult{
		Tool: tools[0], // 取第一个 tool
	}

	fmt.Printf("✅ Planner result: tool=%s\n\n", plannerResult.Tool)

	return plannerResult, nil
}
