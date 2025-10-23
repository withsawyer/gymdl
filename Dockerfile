# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 下载依赖，利用缓存加速
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源代码
COPY . .

# 编译静态二进制
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# Final stage
FROM alpine:latest
WORKDIR /root/

# 只复制编译后的二进制
COPY --from=builder /app/app .

# 容器启动命令
CMD ["./app"]
