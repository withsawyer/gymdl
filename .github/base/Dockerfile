FROM python:3.12-slim

# 设置时区
ENV TZ=Asia/Shanghai \
    DEBIAN_FRONTEND=noninteractive

# 安装构建工具、Go、ffmpeg编译依赖
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    tzdata ca-certificates curl wget git build-essential pkg-config \
    yasm nasm libtool autoconf automake cmake golang-go && \
    ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    rm -rf /var/lib/apt/lists/*

# -------------------------------------------------------------
# 编译安装 libfdk_aac
# -------------------------------------------------------------
WORKDIR /tmp/build
ENV FFMPEG_VERSION=6.1 \
    FDKAAC_VERSION=2.0.2

RUN wget https://github.com/mstorsjo/fdk-aac/archive/refs/tags/v${FDKAAC_VERSION}.tar.gz && \
    tar xzf v${FDKAAC_VERSION}.tar.gz && \
    cd fdk-aac-${FDKAAC_VERSION} && \
    autoreconf -fiv && ./configure --enable-static --disable-shared && \
    make -j$(nproc) && make install

# -------------------------------------------------------------
# 编译安装 ffmpeg (静态 + libfdk_aac)
# -------------------------------------------------------------
RUN wget https://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.gz && \
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
    make -j$(nproc) && make install && \
    strip /usr/local/bin/ffmpeg && strip /usr/local/bin/ffprobe

WORKDIR /root

ENTRYPOINT ["python3"]
