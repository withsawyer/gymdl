package processor

import (
	"fmt"
	"github.com/nichuanfang/gymdl/config"
	"os"
	"path/filepath"
	"time"

	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/utils"
)

// 顶级接口定义

type Processor interface {
	Init(cfg *config.Config)
	Handle(link string) (string, error) // 处理链接并返回结果
	Category() domain.ProcessorCategory // 所属分类
}

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
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("创建目录失败: %v\n", err)
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
		return nil // 不存在则不需要删除
	}

	// 删除整个目录树
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("删除目录失败: %v", err)
	}
	utils.DebugWithFormat("🧹 已删除临时目录: %s\n", dir)
	return nil
}
