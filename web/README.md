<div align="center">
	<img src="../public/favicon.svg" width="160" />
	<h1>ModelGate Web</h1>
</div>

---

[![license](https://img.shields.io/badge/license-MIT-green.svg)](./LICENSE)

> ModelGate 管理后台前端项目 - 一个基于 Vue3、Vite、TypeScript 和 NaiveUI 的现代化 AI 模型服务管理平台

## 简介

ModelGate Web 是 ModelGate 项目的管理后台前端，基于 [SoybeanAdmin](https://github.com/soybeanjs/soybean-admin) 模板开发。它提供了一个清新优雅、高颜值的 Web 界面，用于管理 AI 模型服务、用户账户、使用统计等功能。

## 技术栈

- **框架**: Vue 3.5+ / TypeScript 5.8+
- **构建工具**: Vite 6.3+
- **UI 组件**: Naive UI 2.41+
- **样式**: UnoCSS / SCSS
- **状态管理**: Pinia 3.0+
- **路由**: Vue Router 4.5+ / Elegant Router
- **国际化**: Vue I18n 11.1+
- **通信协议**: gRPC-Web (ConnectRPC)
- **图表**: ECharts 5.6+
- **工具库**: VueUse / Day.js / Clipboard


## 环境要求

- **Node.js**: >= 18.20.0
- **pnpm**: >= 8.7.0

## 快速开始

### 安装依赖

```bash
pnpm install
```

### 开发模式

```bash
# 测试环境
pnpm dev

# 生产环境
pnpm dev:prod
```

默认访问地址: `http://localhost:9527`

### 构建生产

```bash
pnpm build
```

### 代码检查

```bash
pnpm lint
```

### 类型检查

```bash
pnpm typecheck
```

## 环境变量

主要环境变量配置 (`.env` 文件):

| 变量名 | 说明 | 默认值 |
|--------|------|--------|
| `VITE_BASE_URL` | 应用基础路径 | `/` |
| `VITE_APP_TITLE` | 应用标题 | `ModelGate` |
| `VITE_PUBLIC_GRPC_SERVICE_URL` | gRPC 服务地址 | `http://localhost:8889` |
| `VITE_AUTH_ROUTE_MODE` | 路由模式 (static/dynamic) | `dynamic` |
| `VITE_ROUTER_HISTORY_MODE` | 路由模式 (hash/history) | `history` |
| `VITE_STORAGE_PREFIX` | 存储前缀 | `MG_` |

## gRPC 通信

项目使用 ConnectRPC 进行 gRPC-Web 通信，定义了以下服务:

- **AuthService** - 认证服务 (登录、注册、令牌刷新)
- **SystemService** - 系统服务 (用户、角色管理)
- **RelayService** - 中继服务 (模型、定价、服务商管理)

配置文件: `src/grpc.ts`

## Git 提交规范

项目使用 Conventional Commits 规范，可通过以下命令创建符合规范的提交:

```bash
# 中文提交
pnpm commit:zh

# 英文提交
pnpm commit
```

## 浏览器支持

推荐使用最新版 Chrome 浏览器进行开发。

| Chrome | Edge | Firefox | Safari |
|--------|------|---------|--------|
| last 2 versions | last 2 versions | last 2 versions | last 2 versions |

## 相关项目

- [ModelGate](https://github.com/yearnfar/modelgate) - 主项目仓库
- [SoybeanAdmin](https://github.com/soybeanjs/soybean-admin) - 基础模板

## 开源协议

MIT License
