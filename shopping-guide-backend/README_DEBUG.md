# 调试指南

## 1. 环境准备

### Redis 配置
项目使用 `redis-server` Docker 容器：
- **端口**: 6380
- **密码**: redis123
- **容器名**: redis-server

启动 Redis：
```bash
cd shopping-guide-backend
docker-compose up -d
```

验证 Redis 连接：
```bash
docker exec redis-server redis-cli -a redis123 ping
```

### MySQL 配置
确保 MySQL 已启动并创建数据库：
```bash
mysql -u root -p < scripts/init_db.sql
```

## 2. VSCode 调试配置

已创建 `.vscode/launch.json` 配置文件。

### 使用方法：

1. 在 VSCode 中打开项目
2. 按 `F5` 或点击"运行和调试"
3. 选择 "Debug Shopping Guide Backend"
4. 程序将在调试模式下启动

### 设置断点：

在代码中点击行号左侧设置断点，例如：
- `internal/handler/chat_handler.go:26` - Chat 方法入口
- `internal/service/chat_service.go` - Service 层逻辑

## 3. 测试 Chat 接口

### 方式一：使用测试脚本

```bash
cd shopping-guide-backend
./scripts/test_chat.sh
```

### 方式二：使用 curl 命令

**测试健康检查：**
```bash
curl http://localhost:8080/health
```

**测试阻塞式 Chat 接口：**
```bash
curl -X POST http://localhost:8080/api/v1/chat \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-session-001",
    "user_id": "test-user-001",
    "message": "我想买一个性价比高的笔记本电脑",
    "style": "xiaohongshu"
  }'
```

**测试流式 Chat 接口：**
```bash
curl -X POST http://localhost:8080/api/v1/chat/stream \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": "test-session-002",
    "user_id": "test-user-001",
    "message": "推荐一些适合学生的手机",
    "style": "dongyuhui"
  }' \
  -N
```

## 4. 常见问题

### 问题 1: Redis 连接失败

**错误信息：**
```
Failed to init Redis: failed to connect redis: EOF
```

**解决方案：**
1. 检查 Redis 容器是否运行：`docker ps | grep redis-server`
2. 检查端口映射：应该是 `0.0.0.0:6380->6379/tcp`
3. 验证连接：`docker exec redis-server redis-cli -a redis123 ping`

### 问题 2: MySQL 连接失败

**解决方案：**
1. 确保 MySQL 服务运行
2. 检查 `configs/config.dev.yaml` 中的密码配置
3. 确保数据库 `shopping_guide` 已创建

### 问题 3: 路由未注册

**错误信息：**
```
404 page not found
```

**解决方案：**
检查 `internal/router/router.go` 确保 chatHandler 已正确注入。

## 5. 开发流程

1. **启动依赖服务**：
   ```bash
   docker-compose up -d
   ```

2. **按 F5 启动调试**

3. **设置断点**：在关键代码位置设置断点

4. **发送测试请求**：使用 curl 或测试脚本

5. **查看变量**：在调试面板查看变量值和调用栈

6. **单步调试**：使用 F10（单步跳过）、F11（单步进入）进行调试

## 6. 配置文件说明

- `configs/config.yaml` - 默认配置
- `configs/config.dev.yaml` - 开发环境配置（会覆盖默认配置）
- `configs/config.prod.yaml` - 生产环境配置

环境通过 `ENV` 环境变量控制，默认为 `dev`。

