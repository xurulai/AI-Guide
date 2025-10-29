package handler

import (
	"net/http"

	"shopping-guide-backend/internal/model"
	"shopping-guide-backend/internal/service"

	"github.com/gin-gonic/gin"
)

// chatHandler 对话处理器实现
type chatHandler struct {
	chatService service.ChatService
}

// NewChatHandler 创建对话处理器
func NewChatHandler(chatService service.ChatService) ChatHandler {
	return &chatHandler{
		chatService: chatService,
	}
}

// Chat 对话接口（阻塞式）
func (h *chatHandler) Chat(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(model.CodeInvalidParams, err.Error()))
		return
	}

	// 调用Service层处理业务逻辑
	resp, err := h.chatService.Chat(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(model.CodeInternalError, err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}

// ChatStream 对话接口（流式）
func (h *chatHandler) ChatStream(c *gin.Context) {
	var req model.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.NewErrorResponse(model.CodeInvalidParams, err.Error()))
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	// 获取流式响应通道
	stream, err := h.chatService.ChatStream(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.NewErrorResponse(model.CodeInternalError, err.Error()))
		return
	}

	// 发送流式响应
	for chunk := range stream {
		c.SSEvent(chunk.Event, chunk.Data)
		c.Writer.Flush()
	}
}
