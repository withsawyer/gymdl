# ==============================
# 第一阶段：编译 FFmpeg + libfdk_aac
# ==============================
FROM python:3.12-slim AS builder

ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive \
    FFMPEG_VERSION=6.1 \
    FDKAAC_VERSION=2.0.2

# 安装构建依赖
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

# 编译 FFmpeg（静态 + libfdk_aac）
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
# 第二阶段：生产镜像
# ==============================
FROM python:3.12-slim

# 安装运行依赖（删除了 libicu-dev）
RUN apt update && apt install -y --no-install-recommends \
    tzdata ca-certificates wget curl tar unzip libicu-dev && \
    ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo Asia/Shanghai > /etc/timezone && \
    rm -rf /var/lib/apt/lists/*

# 安装 Deno 并清理缓存
RUN curl -fsSL https://deno.land/install.sh | sh && \
    rm -rf /root/.cache/deno 

# ==============================
# 安装 N_m3u8DL-RE
# ==============================
WORKDIR /tmp
RUN curl -L -o N_m3u8DL-RE.tar.gz \
    https://github.com/nilaoda/N_m3u8DL-RE/releases/download/v0.3.0-beta/N_m3u8DL-RE_v0.3.0-beta_linux-x64_20241203.tar.gz && \
    tar -xzf N_m3u8DL-RE.tar.gz && \
    mv N_m3u8DL-RE /usr/local/bin/N_m3u8DL-RE && \
    chmod +x /usr/local/bin/N_m3u8DL-RE && \
    rm -rf /tmp/*

WORKDIR /app
RUN wget -q "https://github.com/zhaarey/wrapper/releases/download/linux.V2/wrapper.x86_64.tar.gz" && \
    wget -q "https://github.com/zhaarey/wrapper/archive/refs/heads/main.zip" &&\
    unzip main.zip &&\
    mv wrapper-main wrapper && \
    tar -xzf wrapper.x86_64.tar.gz -C wrapper && \
    rm -rf wrapper.x86_64.tar.gz && \
    rm -rf main.zip && \
    cd wrapper && \
    chmod +x wrapper

# 拷贝 FFmpeg 静态二进制
COPY --from=builder /usr/local/bin/ffmpeg /usr/local/bin/ffprobe /usr/local/bin/

ENV TZ=Asia/Shanghai \
    PATH="/root/.deno/bin:/app/data/bin:${PATH}" \
    DEBIAN_FRONTEND=noninteractive