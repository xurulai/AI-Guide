package client

import (
	"context"

	"shopping-guide-backend/internal/model"
)

// DifyClient Dify客户端接口
type DifyClient interface {
	CallWorkflow(ctx context.Context, appID string, inputs map[string]interface{}, user string) (*model.DifyWorkflowResponse, error)
	CallWorkflowStream(ctx context.Context, appID string, inputs map[string]interface{}, user string) (<-chan model.StreamChunk, error)
}

// SearchClient 搜推系统客户端接口
type SearchClient interface {
	// TODO: 定义搜推系统接口方法
}

// OrderClient 订单系统客户端接口
type OrderClient interface {
	// TODO: 定义订单系统接口方法
}
