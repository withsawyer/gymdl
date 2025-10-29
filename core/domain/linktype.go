package domain

// LinkType 是所有解析出来的类型枚举
type LinkType string

const (
	/* ---------------------- 音乐平台枚举 ---------------------- */
	LinkUnknown      LinkType = ""
	LinkAppleMusic   LinkType = "AppleMusic"
	LinkNetEase      LinkType = "Netease"
	LinkQQMusic      LinkType = "QQMusic"
	LinkSoundcloud   LinkType = "Soundcloud"
	LinkSpotify      LinkType = "Spotify"
	LinkYoutubeMusic LinkType = "YoutubeMusic"

	/* ---------------------- 视频平台枚举 ---------------------- */

	LinkBilibili    LinkType = "Bilibili"
	LinkDouyin      LinkType = "Douyin"
	LinkXiaohongshu LinkType = "Xiaohongshu"
	LinkYoutube     LinkType = "Youtube"
)
