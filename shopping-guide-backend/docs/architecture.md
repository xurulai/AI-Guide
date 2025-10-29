# 架构说明

## 核心架构设计

### 多Agent架构（一主多从、单步规划）

```
用户请求
    ↓
ChatService (对外接口)
    ↓
OrchestratorService (编排器)
    ├─→ PlannerService (主Agent - 规划器)
    │       ↓
    │   [调用Dify Planner工作流]
    │       ↓
    │   返回Tool选择和ToolInput
    │
    └─→ ExecutorService (从Agent - 执行器)
            ↓
        根据Tool类型路由到对应Executor:
        - PRODUCT_RECOMMENDATION_MODULE (商品推荐Agent)
        - SHOPPING_GUIDE_AND_INTENT_MINING_MODULE (售前导购Agent)  
        - ECOMMERCE_QA_ASSISTANT_MODULE (答疑助手Agent)
            ↓
        [调用对应的Dify Executor工作流]
            ↓
        返回执行结果
```

## 目录结构

```
shopping-guide-backend/
├── cmd/server/main.go           # 程序入口 ✅
├── internal/
│   ├── database/                # 数据库连接初始化 ✅
│   │   ├── redis.go
│   │   └── mysql.go
│   ├── model/                   # 数据模型 ✅
│   ├── config/                  # 配置管理 ✅
│   ├── service/                 # 业务逻辑层
│   │   ├── service.go           # 基础服务接口定义
│   │   ├── chat_service.go      # 对话服务（顶层）
│   │   ├── orchestrator_service.go  # 编排服务 ⭐核心
│   │   ├── planner_service.go   # Planner服务（主Agent）⭐
│   │   ├── executor_service.go  # Executor服务（从Agent）⭐
│   │   ├── session_service.go   # 会话服务
│   │   └── product_service.go   # 商品服务
│   ├── repository/              # 数据访问层接口
│   ├── client/                  # 外部服务客户端接口
│   ├── handler/                 # HTTP处理器
│   ├── middleware/              # 中间件
│   └── router/                  # 路由
├── configs/                     # 配置文件 ✅
├── scripts/                     # 脚本 ✅
├── deploy/                      # 部署文件 ✅
└── docs/                        # 文档 ✅
```

## 核心服务说明

### 1. PlannerService (主Agent - 规划器)

**职责：**
- 分析用户输入的真实购物意图
- 判断用户提及的商品是否在商品库中
- 选择合适的Tool（商品推荐/导购引导/答疑）
- 生成ToolInput供Executor使用

**独立性：**
- 不依赖Executor
- 只负责规划决策，不执行具体业务
- 调用Dify的Planner工作流

### 2. ExecutorService (从Agent - 执行器)

**职责：**
- 根据Planner的Tool选择，路由到对应的Executor
- 调用对应的Dify Executor工作流
- 执行具体的业务逻辑（商品推荐/导购/答疑）

**三个Executor:**
1. **商品推荐Agent**: 检索商品、生成推荐话术
2. **售前导购Agent**: 需求挖掘、引导话术生成
3. **答疑助手Agent**: FAQ检索、答疑话术生成

### 3. OrchestratorService (编排器)

**职责：**
- 协调Planner和Executor的调用流程
- 管理会话上下文
- 处理错误和降级
- 记录日志

**核心流程：**
```go
1. 加载会话 → SessionService
2. 获取商品库 → ProductService
3. 调用Planner规划 → PlannerService
4. 调用Executor执行 → ExecutorService
5. 保存会话和日志
6. 返回响应
```

## Redis使用说明

### 初始化位置
- **位置**: `internal/database/redis.go`
- **调用**: `cmd/server/main.go` 中初始化

### 使用位置
- **Repository层**: 在具体的Repository实现中使用
  - `SessionRepository`: 会话存储（Redis）
  - `CacheRepository`: 缓存存储（Redis）

### 数据结构设计

```
# 会话存储
Key: session:{session_id}
Type: String (JSON)
TTL: 7天

# 用户会话索引
Key: user:sessions:{user_id}
Type: Set
Members: [session_id1, session_id2...]
TTL: 7天

# Dify调用缓存
Key: dify:cache:{md5(inputs)}
Type: String (JSON)
TTL: 1小时

# 限流计数
Key: rate_limit:{user_id}:{endpoint}
Type: String
TTL: 60秒
```

## MySQL使用说明

### 初始化位置
- **位置**: `internal/database/mysql.go`
- **调用**: `cmd/server/main.go` 中初始化

### 使用位置
- **Repository层**: 在具体的Repository实现中使用
  - `ProductRepository`: 商品数据
  - `UserRepository`: 用户数据
  - `LogRepository`: 日志数据

## 依赖注入顺序

```
1. 数据库连接 (MySQL, Redis)
    ↓
2. Repository层
    ↓
3. Client层 (DifyClient)
    ↓
4. Service层 (ProductService, SessionService)
    ↓
5. Service层 (PlannerService, ExecutorService)
    ↓
6. Service层 (OrchestratorService)
    ↓
7. Service层 (ChatService)
    ↓
8. Handler层
    ↓
9. Router
```

## 下一步开发

1. 实现Repository层（连接Redis和MySQL）
2. 实现Client层（DifyClient）
3. 实现Service层业务逻辑
4. 实现Handler层
5. 完善中间件
6. 编写单元测试
