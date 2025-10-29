package model

import "time"

// ChatLog 对话日志
type ChatLog struct {
	LogID               int64                  `json:"log_id" gorm:"primaryKey;autoIncrement;column:log_id"`
	SessionID           string                 `json:"session_id" gorm:"column:session_id;index"`
	UserID              string                 `json:"user_id" gorm:"column:user_id;index"`
	Query               string                 `json:"query" gorm:"column:query;type:text"`
	Response            string                 `json:"response" gorm:"column:response;type:text"`
	ToolUsed            string                 `json:"tool_used" gorm:"column:tool_used;index"`
	PlannerResult       map[string]interface{} `json:"planner_result" gorm:"serializer:json;column:planner_result"`
	ExecutorResult      map[string]interface{} `json:"executor_result" gorm:"serializer:json;column:executor_result"`
	RecommendedProducts []RecommendedProduct   `json:"recommended_products" gorm:"serializer:json;column:recommended_products"`
	LatencyMs           int                    `json:"latency_ms" gorm:"column:latency_ms"`
	TokensUsed          int                    `json:"tokens_used" gorm:"column:tokens_used"`
	CreatedAt           time.Time              `json:"created_at" gorm:"column:created_at;index"`
}

// TableName 指定表名
func (ChatLog) TableName() string {
	return "chat_logs"
}

// DifyCallLog Dify调用日志
type DifyCallLog struct {
	LogID         int64                  `json:"log_id" gorm:"primaryKey;autoIncrement;column:log_id"`
	WorkflowName  string                 `json:"workflow_name" gorm:"column:workflow_name;index"`
	AppID         string                 `json:"app_id" gorm:"column:app_id;index"`
	WorkflowRunID string                 `json:"workflow_run_id" gorm:"column:workflow_run_id"`
	Inputs        map[string]interface{} `json:"inputs" gorm:"serializer:json;column:inputs"`
	Outputs       map[string]interface{} `json:"outputs" gorm:"serializer:json;column:outputs"`
	Status        string                 `json:"status" gorm:"column:status;index"` // success/error/timeout
	ErrorMessage  string                 `json:"error_message" gorm:"column:error_message;type:text"`
	LatencyMs     int                    `json:"latency_ms" gorm:"column:latency_ms"`
	TokensUsed    int                    `json:"tokens_used" gorm:"column:tokens_used"`
	CreatedAt     time.Time              `json:"created_at" gorm:"column:created_at;index"`
}

// TableName 指定表名
func (DifyCallLog) TableName() string {
	return "dify_call_logs"
}

// ProductRecommendation 商品推荐记录
type ProductRecommendation struct {
	RecID      int64     `json:"rec_id" gorm:"primaryKey;autoIncrement;column:rec_id"`
	SessionID  string    `json:"session_id" gorm:"column:session_id;index"`
	UserID     string    `json:"user_id" gorm:"column:user_id;index"`
	ProductID  string    `json:"product_id" gorm:"column:product_id;index"`
	Reason     string    `json:"reason" gorm:"column:reason;type:text"`
	UserAction string    `json:"user_action" gorm:"column:user_action;index"` // view/click/add_cart/purchase
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
}

// TableName 指定表名
func (ProductRecommendation) TableName() string {
	return "product_recommendations"
}
