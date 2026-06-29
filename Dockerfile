# 多阶段构建 - clawwatch-hub
# Stage 1: 编译
FROM --platform=linux/amd64 golang:1.22-alpine AS builder

WORKDIR /build

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 编译
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o clawwatch-hub .

# Stage 2: 运行时（最小镜像）
FROM --platform=linux/amd64 alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=builder /build/clawwatch-hub .

EXPOSE 4848

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget -qO- http://localhost:4848/api/health || exit 1

CMD ["./clawwatch-hub", "--port", "4848"]
