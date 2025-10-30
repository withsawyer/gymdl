package factory

import (
    "github.com/nichuanfang/gymdl/config"
    "github.com/nichuanfang/gymdl/core/domain"
    "github.com/nichuanfang/gymdl/processor"
    "github.com/nichuanfang/gymdl/processor/music"
    "github.com/nichuanfang/gymdl/processor/video"
)

// GetProcessor 处理器工厂
func GetProcessor(linkType domain.LinkType, cfg *config.Config) processor.Processor{
	switch linkType {
	case domain.LinkAppleMusic:
		return music.NewAppleMusicProcessor(cfg, music.AppleMusicTempDir)
	case domain.LinkNetEase:
		return music.NewNetEaseProcessor(cfg, music.NCMTempDir)
	case domain.LinkQQMusic:
		return music.NewQQMusicProcessor(cfg, music.QQTempDir)
	case domain.LinkSoundcloud:
		return music.NewSoundCloudProcessor(cfg, music.SoundcloudTempDir)
	case domain.LinkSpotify:
		return music.NewSpotifyProcessor(cfg, music.SpotifyTempDir)
	case domain.LinkYoutubeMusic:
		return music.NewYoutubeMusicProcessor(cfg, music.YoutubeTempDir)
	case domain.LinkBilibili:
		return video.NewBiliBiliProcessor(cfg, video.BilibiliTempDir)
	case domain.LinkDouyin:
		return video.NewDouyinProcessor(cfg, video.DouyinTempDir)
	case domain.LinkXiaohongshu:
		return video.NewXiaohongshuProcessor(cfg, video.XHSTempDir)
	case domain.LinkYoutube:
		return video.NewYoutubeProcessor(cfg, video.YoutubeTempDir)
	default:
		return nil
	}
}
