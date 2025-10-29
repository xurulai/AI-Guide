package model

// APIResponse 统一API响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 标准响应码
const (
	CodeSuccess            = 0
	CodeInvalidParams      = 400
	CodeUnauthorized       = 401
	CodeForbidden          = 403
	CodeNotFound           = 404
	CodeRateLimited        = 429
	CodeInternalError      = 500
	CodeServiceUnavailable = 503
)

// NewSuccessResponse 创建成功响应
func NewSuccessResponse(data interface{}) *APIResponse {
	return &APIResponse{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	}
}

// NewErrorResponse 创建错误响应
func NewErrorResponse(code int, message string) *APIResponse {
	return &APIResponse{
		Code:    code,
		Message: message,
	}
}
