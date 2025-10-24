# -------------------------------------------------------------
# 基础镜像: 使用 Debian 作为构建基础
# -------------------------------------------------------------
FROM debian:bookworm

# -------------------------------------------------------------
# 设置时区和基础依赖
# -------------------------------------------------------------
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    tzdata ca-certificates curl wget git build-essential pkg-config \
    yasm nasm libtool autoconf automake cmake python3.12 python3.12-venv python3.12-dev \
    golang-go && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    rm -rf /var/lib/apt/lists/*

# -------------------------------------------------------------
# 编译安装 ffmpeg + libfdk_aac
# -------------------------------------------------------------
WORKDIR /tmp/build

# 版本可根据需求调整
ENV FFMPEG_VERSION=6.1 \
    FDKAAC_VERSION=2.0.2

# 下载并编译 libfdk_aac
RUN wget https://github.com/mstorsjo/fdk-aac/archive/refs/tags/v${FDKAAC_VERSION}.tar.gz && \
    tar xzf v${FDKAAC_VERSION}.tar.gz && \
    cd fdk-aac-${FDKAAC_VERSION} && \
    autoreconf -fiv && ./configure --enable-static --disable-shared && \
    make -j$(nproc) && make install

# 下载并编译 ffmpeg (静态 + fdk_aac)
RUN wget https://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz && \
    tar xzf ffmpeg-${FFMPEG_VERSION}.tar.gz && \
    cd ffmpeg-${FFMPEG_VERSION} && \
    ./configure \
      --pkg-config-flags="--static" \
      --extra-cflags="-I/usr/local/include" \
      --extra-ldflags="-L/usr/local/lib" \
      --extra-libs="-
