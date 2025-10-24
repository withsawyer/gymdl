# ==========================
# 构建阶段
# ==========================
FROM python:3.12-alpine AS builder

ENV FFMPEG_VERSION=6.1.2
ENV FDK_AAC_VERSION=2.0.3
ENV PREFIX=/usr/local
ENV BUILD_DIR=/tmp/ffmpeg_build
ENV PATH="$PREFIX/bin:$PATH"

# 安装构建依赖
RUN apk update && apk add --no-cache \
        build-base yasm nasm pkgconf autoconf automake libtool git wget tar coreutils \
        openssl-dev lame-dev libvpx-dev opus-dev libogg-dev libvorbis-dev libwebp-dev zlib-dev

# 创建构建目录
RUN mkdir -p "$BUILD_DIR"

# 编译 libfdk_aac
RUN cd "$BUILD_DIR" \
    && git clone https://github.com/mstorsjo/fdk-aac.git \
    && cd fdk-aac \
    && git checkout v${FDK_AAC_VERSION} \
    && autoreconf -fiv \
    && ./configure --prefix="$PREFIX" --disable-shared --enable-static \
    && make -j"$(nproc)" \
    && make install

# 编译 FFmpeg
RUN cd "$BUILD_DIR" \
    && wget -O ffmpeg-${FFMPEG_VERSION}.tar.bz2 "https://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.bz2" \
    && tar xjf ffmpeg-${FFMPEG_VERSION}.tar.bz2 \
    && cd ffmpeg-${FFMPEG_VERSION} \
    && PKG_CONFIG_PATH="$PREFIX/lib/pkgconfig" ./configure \
        --prefix="$PREFIX" \
        --pkg-config-flags="--static" \
        --extra-cflags="-I$PREFIX/include" \
        --extra-ldflags="-L$PREFIX/lib" \
        --extra-libs="-lpthread -lm" \
        --enable-gpl \
        --enable-nonfree \
        --enable-libfdk-aac \
        --enable-libmp3lame \
        --enable-libvpx \
        --enable-libopus \
        --enable-libvorbis \
        --enable-libwebp \
        --disable-debug \
        --disable-doc \
        --disable-ffplay \
        --enable-static \
        --disable-shared \
    && make -j"$(nproc)" \
    && make install \
    && make distclean

# 清理构建目录
RUN rm -rf "$BUILD_DIR"

# ==========================
# 运行阶段
# ==========================
FROM python:3.12-alpine

ENV PREFIX=/usr/local
ENV PATH="$PREFIX/bin:$PATH"

# 设置时区并安装运行时依赖
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata

# 拷贝 ffmpeg 二进制文件和库
COPY --from=builder /usr/local /usr/local

# 验证 ffmpeg 可用
RUN ffmpeg -version

ENTRYPOINT ["python3"]