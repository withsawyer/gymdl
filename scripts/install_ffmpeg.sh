#!/bin/sh
set -e

#=====================================================
# 配置部分
#=====================================================
FFMPEG_VERSION="6.1.2"
FDK_AAC_VERSION="2.0.3"
PREFIX="/app/data"
BUILD_DIR="/tmp/ffmpeg_build"

#=====================================================
# 安装构建依赖
#=====================================================
apk update
apk add --no-cache \
    build-base \
    yasm \
    nasm \
    pkgconf \
    autoconf \
    automake \
    libtool \
    git \
    wget \
    tar \
    coreutils \
    openssl-dev \
    lame-dev \
    libvpx-dev \
    opus-dev \
    libogg-dev \
    libvorbis-dev \
    libwebp-dev \
    zlib-dev

#=====================================================
# 创建构建目录
#=====================================================
mkdir -p "${BUILD_DIR}"
cd "${BUILD_DIR}"

#=====================================================
# 编译 libfdk_aac
#=====================================================
echo ">>> 编译 libfdk_aac ${FDK_AAC_VERSION} ..."
git clone https://github.com/mstorsjo/fdk-aac.git
cd fdk-aac
git checkout v${FDK_AAC_VERSION}
autoreconf -fiv
./configure --prefix="${PREFIX}" --disable-shared --enable-static
make -j"$(nproc)"
make install
cd ..

#=====================================================
# 编译 FFmpeg
#=====================================================
echo ">>> 编译 FFmpeg ${FFMPEG_VERSION} ..."
wget -O ffmpeg-${FFMPEG_VERSION}.tar.bz2 "https://ffmpeg.org/releases/ffmpeg-${FFMPEG_VERSION}.tar.bz2"
tar xjf ffmpeg-${FFMPEG_VERSION}.tar.bz2
cd ffmpeg-${FFMPEG_VERSION}

PKG_CONFIG_PATH="${PREFIX}/lib/pkgconfig"
export PKG_CONFIG_PATH

./configure \
    --prefix="${PREFIX}" \
    --pkg-config-flags="--static" \
    --extra-cflags="-I${PREFIX}/include" \
    --extra-ldflags="-L${PREFIX}/lib" \
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
    --disable-shared

make -j"$(nproc)"
make install
make distclean

#=====================================================
# 清理与验证
#=====================================================
cd /
rm -rf "${BUILD_DIR}"

echo ">>> FFmpeg 构建完成！"

# 赋予可执行权限
chmod +x "${PREFIX}/bin/ffmpeg" "${PREFIX}/bin/ffprobe"

# 验证版本
"${PREFIX}/bin/ffmpeg" -version
"${PREFIX}/bin/ffprobe" -version