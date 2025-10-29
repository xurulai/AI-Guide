package model

// ChatRequest 对话请求
type ChatRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Query     string `json:"query" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	//Context   map[string]interface{} `json:"context,omitempty"`
}

// ChatResponse 对话响应
type ChatResponse struct {
	SessionID           string               `json:"session_id"`
	Response            string               `json:"response"`
	ToolUsed            string               `json:"tool_used"`
	RecommendedProducts []RecommendedProduct `json:"recommended_products,omitempty"`
	Metadata            ChatMetadata         `json:"metadata"`
}

// ChatMetadata 对话元数据
type ChatMetadata struct {
	PlannerResult PlannerResult `json:"planner_result"`
	LatencyMs     int64         `json:"latency_ms"`
	TokensUsed    int           `json:"tokens_used"`
}

// SessionCreateRequest 创建会话请求
type SessionCreateRequest struct {
	UserID              string   `json:"user_id" binding:"required"`
	BusinessInstruction string   `json:"business_instruction"`
	ProductCategories   []string `json:"product_categories"`
}

// SessionCreateResponse 创建会话响应
type SessionCreateResponse struct {
	SessionID string `json:"session_id"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at"`
}

// SessionListRequest 会话列表请求
type SessionListRequest struct {
	UserID string `form:"user_id" binding:"required"`
	Page   int    `form:"page" binding:"min=1"`
	Size   int    `form:"size" binding:"min=1,max=100"`
}

// SessionListResponse 会话列表响应
type SessionListResponse struct {
	Sessions []SessionSummary `json:"sessions"`
	Total    int              `json:"total"`
	Page     int              `json:"page"`
	Size     int              `json:"size"`
}

// SessionSummary 会话摘要
type SessionSummary struct {
	SessionID      string `json:"session_id"`
	UserID         string `json:"user_id"`
	MessageCount   int    `json:"message_count"`
	LastMessage    string `json:"last_message"`
	LastUpdateTime string `json:"last_update_time"`
}
