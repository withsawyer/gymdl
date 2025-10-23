# 1. 构建阶段
FROM golang:1.24 AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 并预先下载依赖（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制所有代码
COPY . .

# 编译可执行文件（静态编译，避免依赖 glibc）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# 2. 运行阶段
FROM python:3.12-alpine

WORKDIR /app

# 安装 tzdata 并设置时区
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 安装系统依赖（gcc、musl-dev 等可能用于某些 Python 包编译）
RUN apk add --no-cache build-base libffi-dev

# 复制 Python 依赖文件
COPY ./internal/cron/dependency/requirements.txt ./

# 安装 Python 依赖
RUN pip install --no-cache-dir -r requirements.txt

# 拷贝 Go 编译好的二进制文件
COPY --from=builder /app/app ./

# 暴露端口
EXPOSE 8080

# 启动程序
ENTRYPOINT ["./app"]