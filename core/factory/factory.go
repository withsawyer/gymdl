package factory

import (
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core/domain"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/processor/music"
	"github.com/nichuanfang/gymdl/processor/video"
)

// GetProcessor 处理器工厂
func GetProcessor(linkType domain.LinkType, cfg *config.Config) processor.Processor {
	switch linkType {
	case domain.LinkAppleMusic:
		return music.NewAppleMusicProcessor(cfg)
	case domain.LinkNetEase:
		return music.NewNetEaseProcessor(cfg)
	case domain.LinkQQMusic:
		return music.NewQQMusicProcessor(cfg)
	case domain.LinkSoundcloud:
		return music.NewSoundCloudProcessor(cfg)
	case domain.LinkSpotify:
		return music.NewSpotifyProcessor(cfg)
	case domain.LinkYoutubeMusic:
		return music.NewYoutubeMusicProcessor(cfg)
	case domain.LinkBilibili:
		return video.NewBiliBiliProcessor(cfg)
	case domain.LinkDouyin:
		return video.NewDouyinProcessor(cfg)
	case domain.LinkXiaohongshu:
		return video.NewXiaohongshuProcessor(cfg)
	case domain.LinkYoutube:
		return video.NewYoutubeProcessor(cfg)
	default:
		return nil
	}
}
