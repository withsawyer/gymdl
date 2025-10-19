package ffmpeg

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/utils"
)

func InitFfmpeg(config *config.FfmpegConfig) {
	SetFfmpeg(config.FfmpegPath)
	SetFfprobe(config.FfprobePath)
}

func SetFfmpeg(path string) {
	if !utils.FileExist(path) {
		return
	}
	ffmpeg = path
}

func SetFfprobe(path string) {
	if !utils.FileExist(path) {
		return
	}

	ffprobe = path
}
