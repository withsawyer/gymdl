FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装 git 和 ca-certificates
RUN apk add --no-cache git ca-certificates

# 下载依赖，利用缓存
COPY go.mod go.sum ./
RUN go env -w GOPROXY=https://proxy.golang.org,direct \
    && go mod download

# 拷贝源代码
COPY . .

# 编译
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/app .
CMD ["./app"]
