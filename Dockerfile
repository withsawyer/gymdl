# 使用官方的Go语言镜像作为构建环境
FROM golang:1.24.9-alpine AS builder

# 设置工作目录
WORKDIR /app

# 复制Go模块文件并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目源码
COPY . .

# 编译Go项目，生成可执行文件app
RUN go build -o app .

# 使用一个更小的镜像来运行可执行文件
FROM alpine:latest

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制可执行文件到当前镜像
COPY --from=builder /app/app .

# 声明容器运行时监听的端口（根据你的程序实际端口修改）
EXPOSE 8080

# 运行可执行文件
CMD ["./app"]
