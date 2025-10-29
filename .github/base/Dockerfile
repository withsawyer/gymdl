# ==============================
# 第一阶段：编译 FFmpeg + libfdk_aac
# ==============================
FROM python:3.12-slim AS ffmpeg-builder

ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive \
    FFMPEG_VERSION=6.1 \
    FDKAAC_VERSION=2.0.2

RUN apt update && apt install -y --no-install-recommends \
    tzdata ca-certificates curl wget git build-essential pkg-config \
    yasm nasm libtool autoconf automake cmake && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /tmp/build

# 编译 libfdk_aac
RUN wget -q https://github.com/mstorsjo/fdk-aac/archive/refs/tags/v${FDKAAC_VERSION}.tar.gz && \
    tar xzf v${FDKAAC_VERSION}.tar.gz && \
    cd fdk-aac-${FDKAAC_VERSION} && \
    autoreconf -fiv && ./configure --enable-static --disable-shared && \
    make -j"$(nproc)" && make install

# 编译 FFmpeg
RUN wget -q https://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz && \
    tar xzf ffmpeg-${FFMPEG_VERSION}.tar.gz && \
    cd ffmpeg-${FFMPEG_VERSION} && \
    ./configure \
      --pkg-config-flags="--static" \
      --extra-cflags="-I/usr/local/include" \
      --extra-ldflags="-L/usr/local/lib" \
      --extra-libs="-lpthread -lm" \
      --enable-gpl \
      --enable-nonfree \
      --enable-libfdk_aac \
      --disable-shared \
      --enable-static && \
    make -j"$(nproc)" && make install && \
    strip /usr/local/bin/ffmpeg /usr/local/bin/ffprobe && \
    rm -rf /tmp/build


# ==============================
# 第二阶段：编译 Go 程序 apple-music-downloader
# ==============================
FROM golang:1.23.1 AS go-builder

WORKDIR /go/src/app

# 下载项目源码
RUN git clone https://github.com/zhaarey/apple-music-downloader.git . && \
    go mod tidy

# 编译（静态可执行）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o amdl .


# ==============================
# 第三阶段：生产镜像
# ==============================
FROM python:3.12-slim-bookworm

# 安装运行依赖
RUN apt update && apt install -y --no-install-recommends \
    tzdata ca-certificates wget curl tar unzip libicu-dev gpac && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo Asia/Shanghai > /etc/timezone && \
    rm -rf /var/lib/apt/lists/* 

# 安装 Deno
RUN curl -fsSL https://deno.land/install.sh | sh && \
    rm -rf /root/.cache/deno

# 安装 N_m3u8DL-RE
WORKDIR /tmp
RUN curl -L -o N_m3u8DL-RE.tar.gz \
    https://github.com/nilaoda/N_m3u8DL-RE/releases/download/v0.3.0-beta/N_m3u8DL-RE_v0.3.0-beta_linux-x64_20241203.tar.gz && \
    tar -xzf N_m3u8DL-RE.tar.gz && \
    mv N_m3u8DL-RE /usr/local/bin/N_m3u8DL-RE && \
    chmod +x /usr/local/bin/N_m3u8DL-RE && \
    rm -rf /tmp/*

# 安装 wrapper
WORKDIR /app
RUN wget -q "https://github.com/zhaarey/wrapper/releases/download/linux.V2/wrapper.x86_64.tar.gz" && \
    wget -q "https://github.com/zhaarey/wrapper/archive/refs/heads/main.zip" &&\
    unzip main.zip &&\
    mv wrapper-main wrapper && \
    tar -xzf wrapper.x86_64.tar.gz -C wrapper && \
    rm -rf wrapper.x86_64.tar.gz main.zip && \
    cd wrapper && chmod +x wrapper

# 拷贝 ffmpeg 二进制
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /usr/local/bin/

# 拷贝 go 程序
COPY --from=go-builder /go/src/app/amdl /usr/local/bin/

ENV TZ=Asia/Shanghai \
    PATH="/root/.deno/bin:/app/data/bin:${PATH}" \
    DEBIAN_FRONTEND=noninteractive