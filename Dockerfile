# ---------------------------
# 1. 构建 Go 二进制
# ---------------------------
FROM golang:1.24 AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 并预先下载依赖（利用缓存）
COPY go.mod go.sum ./
RUN go mod download

# 复制所有代码
COPY . .

# 编译静态二进制
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# ---------------------------
# 2. 运行阶段
# ---------------------------
FROM python:3.12-alpine

WORKDIR /app

# 安装基础依赖和 tzdata
RUN apk add --no-cache \
        tzdata \
        build-base \
        libffi-dev \
        wget \
        xz \
        tar \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

COPY requirements.txt ./
RUN mkdir data

# 复制 Go 编译好的二进制文件
COPY --from=builder /app/app ./

# 暴露端口
EXPOSE 8080

# 启动程序
ENTRYPOINT ["./app"]
