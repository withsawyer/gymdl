package linkparser

import (
	"regexp"

	"github.com/nichuanfang/gymdl/core/domain"
)

// 处理器规则集
// 处理器规则集
var linkTypeMatchers = []linkTypeMatcher{
	/* ---------------------- 网易云音乐 ---------------------- */
	{
		domains: []string{
			"music.163.com",
			"y.music.163.com",
			"163cn.tv",
			"163cn.link",
		},
		patterns: []*regexp.Regexp{
			// 网页 / 移动端链接（song / playlist / album / artist），id 为数字
			regexp.MustCompile(`^https?://(?:music\.163\.com|y\.music\.163\.com)/(?:song|playlist|album|artist)\?id=\d+(?:&.*)?$`),

			// 网易云短链（App 内分享）
			regexp.MustCompile(`^https?://163cn\.tv/[A-Za-z0-9]+(?:\?.*)?$`),

			// 其他短链（新格式）
			regexp.MustCompile(`^https?://163cn\.link/[A-Za-z0-9]+(?:\?.*)?$`),
		},
		linkType: domain.LinkNetEase,
	},

	/* ---------------------- YouTube / YouTube Music ---------------------- */
	{
		domains: []string{"youtube.com", "music.youtube.com", "youtu.be"},
		patterns: []*regexp.Regexp{
			// 普通 YouTube 视频
			regexp.MustCompile(`^https?://(?:www\.)?youtube\.com/watch\?v=[\w-]+(?:&.*)?$`),

			// YouTube Music 视频
			regexp.MustCompile(`^https?://music\.youtube\.com/watch\?v=[\w-]+(?:&.*)?$`),

			// 短链格式
			regexp.MustCompile(`^https?://youtu\.be/[\w-]+(?:\?.*)?$`),
		},
		linkType: domain.LinkYoutubeMusic,
	},

	/* ---------------------- Apple Music ---------------------- */
	{
		domains: []string{"music.apple.com"},
		patterns: []*regexp.Regexp{
			// 公共播放列表（pl.u- 开头）
			regexp.MustCompile(`^https?://music\.apple\.com/[a-z]{2}/playlist/[A-Za-z0-9._%\-]+/pl\.u-[A-Za-z0-9]+(?:\?.*)?$`),

			// 资料库播放列表（p. 开头）
			regexp.MustCompile(`^https?://music\.apple\.com/library/playlist/p\.[A-Za-z0-9]+(?:\?.*)?$`),

			// 单曲（song）
			regexp.MustCompile(`^https?://music\.apple\.com/[a-z]{2}/song/[A-Za-z0-9%._\-]+/\d+(?:\?.*)?$`),

			// 专辑（album）
			regexp.MustCompile(`^https?://music\.apple\.com/[a-z]{2}/album/[A-Za-z0-9%._\-]+/\d+(?:\?.*)?$`),
		},
		linkType: domain.LinkAppleMusic,
	},

	/* ---------------------- SoundCloud ---------------------- */
	{
		domains: []string{"soundcloud.com", "snd.sc"},
		patterns: []*regexp.Regexp{
			// 用户主页 / 曲目 / 播放列表
			regexp.MustCompile(`^https?://(?:soundcloud\.com|snd\.sc)/[A-Za-z0-9_\-]+/(?:sets/[A-Za-z0-9_\-]+|[A-Za-z0-9_\-]+)(?:\?.*)?$`),
		},
		linkType: domain.LinkSoundcloud,
	},

	/* ---------------------- QQ 音乐 ---------------------- */
	{
		domains: []string{"y.qq.com", "c.y.qq.com", "m.y.qq.com"},
		patterns: []*regexp.Regexp{
			// 支持 song / album / playlist + id 参数 或 URL 路径形式
			regexp.MustCompile(`^https?://(?:y\.qq\.com|c\.y\.qq\.com|m\.y\.qq\.com)/(?:song|album|playlist)(?:/[A-Za-z0-9_\-]+)?(?:\?id=\d+|/[\dA-Za-z]+)(?:&.*)?$`),
		},
		linkType: domain.LinkQQMusic,
	},

	/* ---------------------- Spotify ---------------------- */
	{
		domains: []string{"open.spotify.com", "play.spotify.com"},
		patterns: []*regexp.Regexp{
			// track / album / playlist + ID (通常 22 字符)
			regexp.MustCompile(`^https?://(?:open\.spotify\.com|play\.spotify\.com)/(?:track|album|playlist)/[A-Za-z0-9]+(?:\?.*)?$`),
		},
		linkType: domain.LinkSpotify,
	},
	/* ---------------------- 待补充 ---------------------- */
}
