#!/bin/bash

# 测试 Chat 接口的脚本

# 设置服务器地址
SERVER_URL="http://localhost:8080"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}========================================${NC}"
echo -e "${YELLOW}测试 Shopping Guide Chat API${NC}"
echo -e "${YELLOW}========================================${NC}\n"

# 1. 健康检查
echo -e "${GREEN}1. 测试健康检查接口${NC}"
curl -X GET "${SERVER_URL}/health" \
    -H "Content-Type: application/json" \
    -w "\n\nHTTP Status: %{http_code}\n\n"

# 2. 测试阻塞式对话接口
echo -e "${GREEN}2. 测试阻塞式对话接口 (POST /api/v1/chat)${NC}"
curl -X POST "${SERVER_URL}/api/v1/chat" \
    -H "Content-Type: application/json" \
    -d '{
        "session_id": "test-session-001",
        "user_id": "test-user-001",
        "message": "我想买一个性价比高的笔记本电脑",
        "style": "xiaohongshu"
    }' \
    -w "\n\nHTTP Status: %{http_code}\n\n"

# 3. 测试流式对话接口
echo -e "${GREEN}3. 测试流式对话接口 (POST /api/v1/chat/stream)${NC}"
curl -X POST "${SERVER_URL}/api/v1/chat/stream" \
    -H "Content-Type: application/json" \
    -d '{
        "session_id": "test-session-002",
        "user_id": "test-user-001",
        "message": "推荐一些适合学生的手机",
        "style": "dongyuhui"
    }' \
    -N \
    -w "\n\nHTTP Status: %{http_code}\n\n"

echo -e "\n${YELLOW}========================================${NC}"
echo -e "${YELLOW}测试完成${NC}"
echo -e "${YELLOW}========================================${NC}"

