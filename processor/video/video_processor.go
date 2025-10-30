package video

import (
	"path/filepath"

	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 视频处理接口定义 ---------------------- */
type Processor interface {
	processor.Processor
	// 视频处理器名称
	Name() domain.LinkType
	// 视频元信息列表
	Videos() []*VideoInfo
	// 下载视频
	Download(url string) error
}

/* ---------------------- 视频结构体定义 ---------------------- */
type VideoInfo struct {
	Title    string
	Author   string
	Duration int
}

/* ---------------------- 常量 ---------------------- */
var BaseTempDir = filepath.Join("data", "temp", "video")

// bilibili临时文件夹
var BilibiliTempDir = filepath.Join(BaseTempDir, "Bilibili")

// 抖音临时文件夹
var DouyinTempDir = filepath.Join(BaseTempDir, "Douyin")

// 小红书临时文件夹
var XHSTempDir = filepath.Join(BaseTempDir, "XHS")

// Youtube临时文件夹
var YoutubeTempDir = filepath.Join(BaseTempDir, "Youtube")

/* ---------------------- 视频下载相关业务函数 ---------------------- */
