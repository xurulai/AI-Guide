package router

import (
	"shopping-guide-backend/internal/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(chatHandler handler.ChatHandler) *gin.Engine {
	r := gin.Default()

	// 中间件
	// r.Use(middleware.CORS())
	// r.Use(middleware.Logger())
	// r.Use(middleware.Recovery())

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API路由组
	v1 := r.Group("/api/v1")
	{
		// 对话接口
		if chatHandler != nil {
			v1.POST("/chat", chatHandler.Chat)
			v1.POST("/chat/stream", chatHandler.ChatStream)
		}

		// 会话接口
		v1.POST("/sessions", func(c *gin.Context) {
			// TODO: sessionHandler.CreateSession
		})
		v1.GET("/sessions/:session_id", func(c *gin.Context) {
			// TODO: sessionHandler.GetSession
		})
		v1.GET("/sessions", func(c *gin.Context) {
			// TODO: sessionHandler.ListSessions
		})
		v1.DELETE("/sessions/:session_id", func(c *gin.Context) {
			// TODO: sessionHandler.DeleteSession
		})
	}

	// 内部接口（供Dify调用）
	internal := r.Group("/internal")
	{
		internal.POST("/products/search", func(c *gin.Context) {
			// TODO: productHandler.SearchProducts
		})
	}

	// 管理接口
	admin := r.Group("/admin")
	{
		admin.GET("/dify/workflows/status", func(c *gin.Context) {
			// TODO: adminHandler.DifyWorkflowStatus
		})
	}

	return r
}
