# ModelGate Dockerfile
# Multi-stage build for minimal image size with Web UI included

# =============================================================================
# Stage 1: Build Web UI
# =============================================================================
FROM node:20-alpine AS web-builder

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

WORKDIR /app/web

# Copy package files
COPY web/pnpm-workspace.yaml web/package.json web/pnpm-lock.yaml ./

# Copy packages
COPY web/packages ./packages

# Install dependencies
RUN pnpm install --frozen-lockfile

# Copy web source code
COPY web/ .

# Build the web UI
RUN pnpm run build

# =============================================================================
# Stage 2: Build Backend
# =============================================================================
FROM golang:1.25.5-alpine AS go-builder

ARG GIT_VERSION
ARG GIT_COMMIT
ARG BUILD_TIME

# Install build dependencies
RUN apk add --no-cache git make

# Install buf
RUN apk add --no-cache curl \
  && curl -sSL -o /usr/local/bin/buf https://github.com/bufbuild/buf/releases/download/v1.63.0/buf-Linux-x86_64 \
  && chmod +x /usr/local/bin/buf

# Set working directory
WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a \
  -installsuffix cgo \
  -ldflags "-X 'main.Version=${GIT_VERSION}' -X 'main.GitCommit=${GIT_COMMIT}' -X 'main.BuildTime=${BUILD_TIME}'" \
  -o modelgate ./cmd/main.go

# =============================================================================
# Stage 3: Runtime
# =============================================================================
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 modelgate && \
  adduser -u 1000 -G modelgate -s /bin/sh -D modelgate

# Set timezone
ENV TZ=Asia/Shanghai

WORKDIR /app

# Copy binary from go-builder
COPY --from=go-builder /build/modelgate .

# Copy configs
COPY --from=go-builder /build/configs ./configs

COPY --from=go-builder /build/migration ./migration

# Copy web UI from web-builder
COPY --from=web-builder /app/web/dist /app/web/dist

# Create directories
RUN mkdir -p logs data && \
  chown -R modelgate:modelgate /app

# Switch to non-root user
USER modelgate

# Expose ports
EXPOSE 8080 8088

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8088/health || exit 1

# Run the application
CMD ["/app/modelgate", "all"]
