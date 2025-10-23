# 1. 构建阶段
FROM golang:1.24 AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制所有代码
COPY . .

# 编译可执行文件
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# 2. 运行阶段
FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/app .

# 暴露端口 (如果你的程序需要)
EXPOSE 8080

ENTRYPOINT ["./app"]
