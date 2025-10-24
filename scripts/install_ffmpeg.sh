#!/bin/sh

# ==============================================================================
# FFmpeg 及 libfdk_aac 编解码器手动安装脚本
#
# 使用方法:
# 1. 启动容器。
# 2. 通过 'docker exec -it [container_id] /app/install_ffmpeg.sh' 手动执行。
# ==============================================================================

echo "Starting FFmpeg and libfdk_aac installation..."

# 1. 启用 Alpine community repository
# FDK AAC 编解码器 (fdk-aac) 位于 Alpine 的 community 仓库
echo "Enabling Alpine Community Repository..."
# 使用临时文件避免直接修改 /etc/apk/repositories 权限问题
TEMP_REPO="/tmp/community.repo"
echo "https://dl-cdn.alpinelinux.org/alpine/v$(awk -F. '{print $1"."$2}' /etc/alpine-release)/community" > $TEMP_REPO
cat $TEMP_REPO >> /etc/apk/repositories
rm $TEMP_REPO

# 2. 更新包索引
apk update

# 3. 安装 ffmpeg 和 fdk-aac (提供了 libfdk_aac)
echo "Installing ffmpeg and fdk-aac..."
apk add --no-cache ffmpeg fdk-aac

# 4. 清理缓存
echo "Cleaning up APK cache..."
rm -rf /var/cache/apk/*

echo "FFmpeg installation complete. You can now run 'ffmpeg -encoders' to verify fdk_aac is available."
