package processor

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/nichuanfang/gymdl/config"

	"github.com/nichuanfang/gymdl/utils"
)

// 顶级接口定义

type Processor interface {
	//构造方法
	Init(cfg *config.Config)
	//处理器名称
	Name() LinkType
}

/* ---------------------- 常量 ---------------------- */

// LinkType 是所有解析出来的类型枚举
type LinkType string

const (
	/* ------------------------ 音乐平台枚举 ---------------------- */
	LinkUnknown      LinkType = ""
	LinkAppleMusic   LinkType = "AppleMusic"
	LinkNetEase      LinkType = "网易云音乐"
	LinkQQMusic      LinkType = "QQ音乐"
	LinkSoundcloud   LinkType = "Soundcloud"
	LinkSpotify      LinkType = "Spotify"
	LinkYoutubeMusic LinkType = "YoutubeMusic"

	/* -------------------------视频平台枚举 ---------------------- */

	LinkBilibili    LinkType = "B站"
	LinkDouyin      LinkType = "抖音"
	LinkXiaohongshu LinkType = "小红书"
	LinkYoutube     LinkType = "Youtube"
)

/* ---------------------- 通用业务工具 ---------------------- */

// BuildOutputDir 构建输出目录
// 规则: baseTempDir + 时间戳（例如：temp/20251030153045）
func BuildOutputDir(baseTempDir string) string {
	// 1. 获取当前时间戳（格式：YYYYMMDDHHMMSS）
	timestamp := time.Now().Format("20060102150405")
	// 2. 构建输出目录路径
	return filepath.Join(baseTempDir, timestamp)
}

// CreateOutputDir 创建临时目录
func CreateOutputDir(outputDir string) error {
	// 判断目录是否存在
	if _, err := os.Stat(outputDir); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		// 其他错误
		return fmt.Errorf("检查目录失败: %v", err)
	}

	// 目录不存在，创建
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}
	utils.DebugWithFormat("🧹 已创建临时目录: %s\n", outputDir)
	return nil
}

// RemoveTempDir 用于清理临时目录
func RemoveTempDir(dir string) error {
	if dir == "" {
		return fmt.Errorf("目录路径为空")
	}

	// 判断目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("检查目录失败: %v", err)
	}

	// 删除整个目录树
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("删除目录失败: %v", err)
	}
	utils.DebugWithFormat("🧹 已删除临时目录: %s\n", dir)
	return nil
}

// DetermineTidyType 获取整理类型
func DetermineTidyType(cfg *config.Config) string {
	return map[int]string{1: "LOCAL", 2: "WEBDAV"}[cfg.Tidy.Mode]
}
