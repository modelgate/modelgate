# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ModelGate is an enterprise-level LLM API gateway and management platform that provides a unified interface to manage and forward requests to multiple LLM services (OpenAI, Anthropic, DeepSeek, Zhipu, etc.). It features user management, API key management, billing system, request logging, and streaming response support.

The project uses a frontend-backend separation architecture:
- **Backend** (Go): API relay service and admin backend service
- **Web** (`/web`): Vue 3 + TypeScript + NaiveUI admin interface

## Ignored Directories
Do not read, analyze, or modify files in the following directories:
- web/
- bin/

## Common Commands

### Backend (Go)

```bash
# Install dependencies
go mod download

# Run all services (API + Admin)
go run cmd/main.go all

# Run API relay service only (port 8888)
go run cmd/main.go api

# Run admin service only (port 8889)
go run cmd/main.go admin

# Database migration
go run cmd/main.go migrate

# Build binary
make build

# Run tests
go test ./...
go test -v ./...
go test -cover ./...

# Run single test
go test -v -run TestFunctionName ./path/to/package

# Lint
make lint

# Format code
go fmt ./...

# Generate protobuf code
buf generate
```

### Frontend (Web)

```bash
cd web

# Install dependencies
pnpm install

# Development server (port 9527)
pnpm dev

# Build for production
pnpm build

# Lint
pnpm lint

# Type check
pnpm typecheck
```

## Architecture

### Backend Structure

```
internal/
├── app/           # Application services (gRPC/HTTP handlers)
│   ├── admin/v1/  # Admin API handlers (auth, system, relay management)
│   └── api/v1/    # Public API handlers (relay endpoint)
├── config/        # Configuration management (Viper, environment variables)
├── relay/         # Relay domain (models, DAOs, services for LLM proxy)
│   ├── dao/       # Data access objects
│   ├── model/     # Domain models (Account, Provider, Model, Ledger, etc.)
│   └── service/   # Business logic for relay operations
├── runtime/       # Core runtime for request processing
│   ├── core/      # Handler, Executor, Hook, Stream interfaces
│   ├── hooks/     # Pre/post request hooks (billing, token counting, etc.)
│   └── provider/  # Provider implementations (openai, anthropic, zhipu, minimax)
├── server/        # HTTP server setup (Gin for API, H2C for Admin gRPC-Web)
│   ├── middleware/ # Auth, CORS, rate limiting, API key validation
│   └── interceptor/ # Logging interceptors
└── system/        # System domain (users, roles, permissions, menus)
    ├── dao/       # Data access objects
    ├── model/     # Domain models (User, Role, Permission, Menu, etc.)
    └── service/   # Business logic for system operations
```

### Dependency Injection

Uses `samber/do/v2` for DI. Services are registered in `cmd/main.go` via `Init()` functions:
- `systemimpl.Init(container)` - System services
- `relayimpl.Init(container)` - Relay services
- `runtime.Init(container)` - Runtime components
- `app.Init(container)` - Application handlers

### Runtime Execution Flow

The relay runtime uses a Handler-Hook-Executor pattern:

1. **Handler** (`internal/runtime/core/handler.go`): Interface for providers with methods:
   - `BeforeRequest`, `DoRequest`, `AfterResponse` (non-streaming)
   - `DoStream` (streaming)

2. **Executor** (`internal/runtime/core/executor.go`): Orchestrates execution:
   - Calls hooks in order (before)
   - Executes handler methods
   - Calls hooks in reverse order (after)
   - Supports retry via `NewRetryExecutor`

3. **Hook** (`internal/runtime/core/hook.go`): Pre/post processing hooks:
   - Billing, token counting, request logging, stream writing

4. **Registry** (`internal/runtime/core/registry.go`): Provider registration:
   - `ExecutorRegistry.Register(provider, NewExecutorFunc)`
   - Providers: openai, anthropic, zhipu, minimax, gemini

### Provider Implementation

Each provider in `internal/runtime/provider/{provider}/` implements:
- Handler with provider-specific request/response transformation
- Model mapping (provider model codes to internal models)
- Client initialization

### Two Server Types

1. **API Server** (port 8888): Gin-based REST API for LLM relay
   - OpenAI-compatible endpoints (`/v1/chat/completions`, `/v1/models`)
   - Middleware: API key auth, rate limiting

2. **Admin Server** (port 8889): H2C (HTTP/2 cleartext) for gRPC-Web
   - Connect RPC handlers for admin operations
   - JWT authentication

### Database

Uses GORM with MySQL. Models defined in:
- `internal/relay/model/` - Relay domain models
- `internal/system/model/` - System domain models

Migration via `go run cmd/main.go migrate` (AutoMigrate).

### Protobuf

Proto files in `/proto`. Generate with `buf generate`:
- Go: `protoc-gen-go`, `protoc-gen-connect-go`
- TypeScript: `protoc-gen-es` (for web client)

## Ports

| Service | Port | Purpose |
|---------|------|---------|
| API Relay | 8888 | OpenAI-compatible API |
| Admin | 8889 | gRPC-Web management API |
| Web Dev | 9527 | Frontend development server |

## Default Credentials

- Username: admin
- Password: 123456

## Environment Configuration

- `.env` - Environment variables (database, redis, secrets)
- `configs/config.toml` - Application configuration