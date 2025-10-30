package video

import (
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
)

/* ---------------------- 视频处理接口定义 ---------------------- */
type VideoProcessor interface {
	processor.Processor
    // 视频处理器名称
	Name() domain.LinkType
    // 视频元信息列表
    Videos() []*VideoInfo
}

/* ---------------------- 视频结构体定义 ---------------------- */
type VideoInfo struct {
	Title    string
	Author   string
	Duration int
}

/* ---------------------- 常量 ---------------------- */

/* ---------------------- 业务工具 ---------------------- */
