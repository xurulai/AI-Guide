package handler

import "github.com/gin-gonic/gin"

// ChatHandler 对话处理器接口
type ChatHandler interface {
	Chat(c *gin.Context)
	ChatStream(c *gin.Context)
}

// SessionHandler 会话处理器接口
type SessionHandler interface {
	CreateSession(c *gin.Context)
	GetSession(c *gin.Context)
	ListSessions(c *gin.Context)
	DeleteSession(c *gin.Context)
}

// ProductHandler 商品处理器接口
type ProductHandler interface {
	SearchProducts(c *gin.Context)
}

// AdminHandler 管理处理器接口
type AdminHandler interface {
	Health(c *gin.Context)
	Metrics(c *gin.Context)
}
