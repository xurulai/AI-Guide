# Handler层最佳实践

## Handler写法对比

### ❌ 不推荐的写法

```go
package handler

import "github.com/gin-gonic/gin"

// 问题1：独立函数，硬编码依赖
func Chat(c *gin.Context) error {  // 问题2：返回error但Gin会忽略
    var req model.ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // 问题3：业务码作为HTTP状态码
        c.JSON(model.CodeInvalidParams, model.NewErrorResponse(model.CodeInvalidParams, err.Error()))
        return err
    }

    // 问题4：直接调用全局service，难以测试
    response, err := service.ChatService(c.Request.Context(), &req)
    if err != nil {
        c.JSON(model.CodeInternalError, model.NewErrorResponse(model.CodeInternalError, err.Error()))
        return err
    }

    c.JSON(model.CodeSuccess, response)
    return nil
}
```

### ✅ 推荐的写法

```go
package handler

import (
    "net/http"
    
    "shopping-guide-backend/internal/model"
    "shopping-guide-backend/internal/service"
    
    "github.com/gin-gonic/gin"
)

// 定义接口（便于mock测试）
type ChatHandler interface {
    Chat(c *gin.Context)
    ChatStream(c *gin.Context)
}

// 实现结构体
type chatHandler struct {
    chatService service.ChatService  // 依赖注入
}

// 构造函数（依赖注入入口）
func NewChatHandler(chatService service.ChatService) ChatHandler {
    return &chatHandler{
        chatService: chatService,
    }
}

// Handler方法
func (h *chatHandler) Chat(c *gin.Context) {
    // 1. 参数验证
    var req model.ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, model.NewErrorResponse(
            model.CodeInvalidParams, 
            err.Error(),
        ))
        return
    }

    // 2. 调用Service
    resp, err := h.chatService.Chat(c.Request.Context(), &req)
    if err != nil {
        // 根据错误类型返回不同的HTTP状态码
        statusCode := http.StatusInternalServerError
        businessCode := model.CodeInternalError
        
        // 可以根据err类型细化
        // if errors.Is(err, errs.ErrSessionNotFound) {
        //     statusCode = http.StatusNotFound
        //     businessCode = model.CodeNotFound
        // }
        
        c.JSON(statusCode, model.NewErrorResponse(businessCode, err.Error()))
        return
    }

    // 3. 返回成功响应
    c.JSON(http.StatusOK, model.NewSuccessResponse(resp))
}
```

## HTTP状态码 vs 业务码

### 分层设计

```
┌─────────────────────────────────────┐
│  HTTP状态码（c.JSON的第一个参数）      │
│  - 200: 成功                         │
│  - 400: 请求参数错误                  │
│  - 401: 未授权                       │
│  - 404: 资源不存在                    │
│  - 500: 服务器内部错误                │
└─────────────────────────────────────┘
              ↓
┌─────────────────────────────────────┐
│  业务码（响应体中的code字段）          │
│  - 0: 成功                           │
│  - 400: 参数错误                     │
│  - 401: 未授权                       │
│  - 404: 资源不存在                    │
│  - 500: 内部错误                     │
└─────────────────────────────────────┘
```

### 响应示例

```json
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "code": 400,              // 业务码
  "message": "invalid parameters",
  "data": null
}
```

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 0,                // 业务码（0表示成功）
  "message": "success",
  "data": {
    "session_id": "xxx",
    "response": "..."
  }
}
```

## 依赖注入示例

### main.go中的组装

```go
func main() {
    // 1. 初始化依赖
    db, _ := database.InitMySQL(&cfg.MySQL)
    rdb, _ := database.InitRedis(&cfg.Redis)
    
    // 2. Repository层
    sessionRepo := repository.NewSessionRepository(rdb, &cfg.Redis)
    productRepo := repository.NewProductRepository(db)
    
    // 3. Client层
    difyClient := client.NewDifyClient(&cfg.Dify)
    
    // 4. Service层（注意依赖顺序）
    productService := service.NewProductService(productRepo)
    sessionService := service.NewSessionService(sessionRepo, productRepo)
    plannerService := service.NewPlannerService(difyClient)
    executorService := service.NewExecutorService(difyClient, productService)
    orchestratorService := service.NewOrchestratorService(
        plannerService,
        executorService,
        sessionService,
        productService,
    )
    chatService := service.NewChatService(orchestratorService)
    
    // 5. Handler层
    chatHandler := handler.NewChatHandler(chatService)
    sessionHandler := handler.NewSessionHandler(sessionService)
    
    // 6. 路由
    r := router.SetupRouter(chatHandler, sessionHandler, cfg)
    r.Run(":8080")
}
```

## 测试示例

### 使用依赖注入的好处：易于测试

```go
package handler

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock ChatService
type MockChatService struct {
    mock.Mock
}

func (m *MockChatService) Chat(ctx context.Context, req *model.ChatRequest) (*model.ChatResponse, error) {
    args := m.Called(ctx, req)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*model.ChatResponse), args.Error(1)
}

// 测试用例
func TestChatHandler_Chat(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    t.Run("success", func(t *testing.T) {
        // 1. 准备mock service
        mockService := new(MockChatService)
        mockService.On("Chat", mock.Anything, mock.Anything).Return(
            &model.ChatResponse{Response: "test response"},
            nil,
        )
        
        // 2. 创建handler（依赖注入mock）
        handler := NewChatHandler(mockService)
        
        // 3. 准备请求
        reqBody := model.ChatRequest{
            SessionID: "test-session",
            Query: "test query",
            UserID: "test-user",
        }
        body, _ := json.Marshal(reqBody)
        
        w := httptest.NewRecorder()
        c, _ := gin.CreateTestContext(w)
        c.Request = httptest.NewRequest("POST", "/chat", bytes.NewBuffer(body))
        c.Request.Header.Set("Content-Type", "application/json")
        
        // 4. 执行
        handler.Chat(c)
        
        // 5. 断言
        assert.Equal(t, http.StatusOK, w.Code)
        mockService.AssertExpectations(t)
    })
}
```

## 总结

| 方面 | 不推荐写法 | 推荐写法 |
|-----|----------|---------|
| **函数类型** | 独立函数 | 结构体方法 |
| **依赖管理** | 硬编码/全局变量 | 依赖注入 |
| **返回值** | `error`（被忽略） | 无返回值 |
| **HTTP状态码** | 混用业务码 | 独立使用 |
| **响应格式** | 不统一 | 统一包装 |
| **可测试性** | 难以测试 | 易于mock |
| **可维护性** | 低 | 高 |

