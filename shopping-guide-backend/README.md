# 导购Agent后端服务

> 电商导购多Agent系统 - Go + Gin + Dify架构

## 项目概述

基于多Agent架构（一主多从、单步规划模式）的智能导购系统，通过Dify平台提供AI能力，Go后端提供业务编排和数据管理。

## 技术栈

- **Web框架**: Gin
- **AI平台**: Dify
- **缓存**: Redis
- **数据库**: MySQL
- **向量库**: Milvus
- **配置**: Viper
- **日志**: Zap
- **HTTP客户端**: Resty

## 目录结构

```
shopping-guide-backend/
├── cmd/                    # 程序入口
├── internal/              # 内部包（不对外暴露）
│   ├── handler/          # HTTP处理器
│   ├── service/          # 业务逻辑层
│   ├── client/           # 外部服务客户端
│   ├── repository/       # 数据访问层
│   ├── model/            # 数据模型
│   ├── middleware/       # 中间件
│   ├── config/           # 配置管理
│   └── pkg/              # 内部工具包
├── configs/              # 配置文件
├── deploy/               # 部署相关
├── docs/                 # 文档
└── scripts/              # 脚本
```

## 快速开始

### 前置要求

- Go 1.21+
- Redis 7+
- MySQL 8+
- Docker & Docker Compose (可选)

### 本地开发

1. 克隆项目
```bash
git clone <repository-url>
cd shopping-guide-backend
```

2. 安装依赖
```bash
go mod download
```

3. 配置环境变量
```bash
cp configs/config.dev.yaml configs/config.local.yaml
# 编辑 config.local.yaml，填入必要的配置
```

4. 初始化数据库
```bash
mysql -u root -p < scripts/init_db.sql
```

5. 启动服务
```bash
go run cmd/server/main.go
```

### Docker部署

```bash
docker-compose up -d
```

## API文档

启动服务后访问：`http://localhost:8080/swagger/index.html`

详细API文档见：[docs/api.md](docs/api.md)

## 配置说明

配置文件位于 `configs/` 目录：
- `config.yaml` - 默认配置
- `config.dev.yaml` - 开发环境
- `config.prod.yaml` - 生产环境

支持通过环境变量覆盖配置，优先级：环境变量 > 本地配置 > 默认配置

## 开发规范

- 遵循Go标准项目布局
- 使用golangci-lint进行代码检查
- 编写单元测试，覆盖率>80%
- 提交前运行 `make test` 和 `make lint`

## 许可

[MIT License](LICENSE)

