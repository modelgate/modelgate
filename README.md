<div align="center">
	<h1>ModelGate</h1>
	<p>一个企业级大模型 API 网关与管理平台</p>
</div>

---

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue)](https://go.dev/)
[![Vue Version](https://img.shields.io/badge/Vue-3.5-brightgreen)](https://vuejs.org/)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

> ModelGate 是一个功能完整的企业级大模型 API 网关，提供统一的接口来管理和转发到多个大模型服务（OpenAI、Anthropic、DeepSeek、智谱等）。支持用户管理、API 密钥管理、计费系统、请求记录等完整功能。

---

## 项目简介

ModelGate 采用前后端分离架构，由以下两部分组成：

- **后端服务** (`/`) - 基于 Go 语言开发，包含 API 转发服务和管理后台服务
- **Web 管理后台** (`/web`) - 基于 Vue 3 + TypeScript + NaiveUI 开发的现代化管理界面

### 核心功能

- **多模型支持** - 支持 OpenAI、Anthropic、智谱 (Zhipu AI) 等多个大模型提供商
- **接口转发** - 提供统一的 API 接口，智能转发到不同的大模型服务
- **用户管理** - 完整的用户注册、登录、认证授权系统
- **API 密钥管理** - 用户 API 密钥的创建、管理和权限控制
- **供应商管理** - 管理不同的大模型供应商及其配置
- **模型管理** - 管理不同供应商的 AI 模型配置和定价
- **计费系统** - 基于 Token 使用量的计费和账户流水管理
- **请求记录** - 完整的 API 请求历史记录和查询
- **流式响应** - 支持流式和非流式响应模式
- **速率限制** - 内置请求速率限制功能
- **Redis 缓存** - 支持 Redis 缓存提升性能

---

## 技术栈

### 后端

| 技术 | 说明 |
|------|------|
| Go 1.25.5 | 开发语言 |
| Gin | REST API 框架 |
| Connect RPC | gRPC-Web 框架 |
| GORM | 数据库 ORM |
| JWT | 认证授权 |
| samber/do | 依赖注入 |
| Viper | 配置管理 |
| Logrus | 日志 |
| Redis | 缓存 |
| Cobra | 命令行工具 |
| Protocol Buffers | 接口定义 |

### 前端

| 技术 | 说明 |
|------|------|
| Vue 3.5+ | 前端框架 |
| TypeScript 5.8+ | 类型系统 |
| Vite 6.3+ | 构建工具 |
| Naive UI 2.41+ | UI 组件库 |
| UnoCSS | 原子化 CSS |
| Pinia 3.0+ | 状态管理 |
| Vue Router 4.5+ | 路由管理 |
| Vue I18n 11.1+ | 国际化 |
| ConnectRPC | gRPC-Web 客户端 |
| ECharts 5.6+ | 图表库 |


## 快速开始

### 环境要求

**后端:**
- Go 1.25.5+
- MySQL 5.7+
- Redis 6.0+
- (可选) Docker

**前端:**
- Node.js 18.20.0+
- pnpm 8.7.0+

### 后端服务

#### 1. 安装依赖

```bash
# 克隆仓库
git clone https://github.com/modelgate/modelgate.git
cd modelgate

# 安装 Go 依赖
go mod download

# 安装 buf (用于 Protobuf 代码生成)
go install github.com/bufbuild/buf/cmd/buf@latest

# 生成 Protobuf 代码
buf generate
```

#### 2. 配置

复制并编辑环境变量：

```bash
cp .env.example .env
```

编辑 `.env` 和 `configs/config.toml`：

```toml
[apiServer]
port = 8888

[adminServer]
port = 8889

[database]
type = "mysql"
host = "127.0.0.1"
port = 3306
name = "modelgate"
user = "your_user"
password = "your_password"

[redis]
host = "localhost"
port = 6379
```

#### 3. 数据库迁移

```bash
go run cmd/main.go migrate
```

#### 4. 启动服务

```bash
# 启动 API 转发服务 (端口 8888)
go run cmd/main.go api

# 启动管理后台服务 (端口 8889)
go run cmd/main.go admin

# 同时启动两个服务
go run cmd/main.go all
```

### 前端项目

#### 1. 安装依赖

```bash
cd web
pnpm install
```

#### 2. 配置

编辑 `web/.env`，设置后端服务地址：

```env
VITE_PUBLIC_GRPC_SERVICE_URL="http://localhost:8889"
```

#### 3. 启动开发服务器

```bash
pnpm dev
```

访问 `http://localhost:9527`

#### 4. 构建生产版本

```bash
pnpm build
```

---

## Docker 部署

### 使用 Docker Compose

```bash
# 启动所有服务
docker-compose -f deployments/docker-compose.yaml up -d
```

### 手动构建

```bash
# 构建镜像
docker build -t modelgate:latest .

# 运行容器
docker run -d \
  -p 8888:8888 \
  -p 8889:8889 \
  -e TD_DATABASE_HOST=your_db_host \
  -e TD_DATABASE_PASSWORD=your_db_password \
  modelgate:latest
```

---

## 开发指南

### 后端开发

#### 开发热重载

使用 [gowatch](https://github.com/silenceper/gowatch)：

```bash
gowatch
```

#### Protobuf 开发

```bash
# 更新 Protobuf 定义后重新生成代码
buf generate

# 验证 Protobuf 文件
buf lint
```

#### 运行测试

```bash
go test ./...
go test -cover ./...
```

### 前端开发

```bash
cd web

# 代码检查
pnpm lint

# 类型检查
pnpm typecheck

# Git 提交 (遵循 Conventional Commits)
pnpm commit:zh
```

---

## 命令行接口

```bash
# 查看帮助
go run cmd/main.go --help

# 启动 API 转发服务
go run cmd/main.go api

# 启动管理后台
go run cmd/main.go admin

# 同时启动两个服务
go run cmd/main.go all

# 数据库迁移
go run cmd/main.go migrate
```

---

## 服务端口

| 服务 | 端口 | 说明 |
|------|------|------|
| API 转发服务 | 8888 | OpenAI 兼容的 API 接口 |
| 管理后台服务 | 8889 | gRPC-Web 管理接口 |
| Web 前端 | 9527 | 管理后台界面 |


## 浏览器支持

推荐使用最新版 Chrome 浏览器进行开发。

| Chrome | Edge | Firefox | Safari |
|--------|------|---------|--------|
| last 2 versions | last 2 versions | last 2 versions | last 2 versions |

---

## 许可证

MIT License - 详见 [LICENSE](LICENSE) 文件

---

## 贡献

欢迎提交 Issue 和 Pull Request！

---

## 联系方式

如有问题或建议，请提交 Issue。

---

<div align="center">
Made with ❤️ by ModelGate Team
</div>
