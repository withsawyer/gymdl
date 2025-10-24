# ---------------------------
# 1. 构建 Go 二进制 (构建阶段)
# ---------------------------
FROM golang:1.24 AS builder

WORKDIR /app

# 复制 go.mod 和 go.sum 并预先下载依赖（利用缓存）
# 这利用了 Docker 缓存，如果依赖未变，则不需要重新下载
COPY go.mod go.sum ./
RUN go mod download

# 复制所有代码
COPY . .

# 编译静态二进制
# CGO_ENABLED=0 确保静态链接，从而避免运行时需要安装额外的 C 库
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

# ---------------------------
# 2. 运行阶段 (纯 Go/Python 运行时环境)
# ---------------------------
FROM ghcr.io/nichuanfang/gymdl-base

WORKDIR /app

# 1. 复制需要的文件
COPY requirements.txt ./

# 2. 复制 Go 编译好的二进制文件
COPY --from=builder /app/app ./

# 3. 把 /app/data/bin 加入 PATH
ENV PATH="/app/data/bin:${PATH}"

# 暴露端口
EXPOSE 8080

# 启动程序
ENTRYPOINT ["./app"]