package downloader

import (
	"github.com/nichuanfang/gymdl/core/downloader/music"
	"github.com/nichuanfang/gymdl/core/downloader/video"
	"net/url"
	"regexp"
	"strings"
	"unicode"
)

// platformMatcher 平台匹配规则
type platformMatcher struct {
	patterns []*regexp.Regexp
	handler  any
	domains  []string // 新增：快速判定域名
}

var platforms = []platformMatcher{
	{
		domains: []string{
			"music.163.com",
			"y.music.163.com",
			"163cn.tv",
			"163cn.link",
		},
		patterns: []*regexp.Regexp{
			// 网页或移动端链接
			regexp.MustCompile(`^https?://(music\.163\.com|y\.music\.163\.com)/(#)?/?(song|playlist|album|artist)\?id=\d+`),
			// 网易云短链（App 内分享）
			regexp.MustCompile(`^https?://163cn\.tv/\w+`),
			// 其他平台
			regexp.MustCompile(`^https?://163cn\.link/\w+`),
		},
		handler: &music.NCMHandler{},
	},
	{
		domains: []string{"youtube.com", "music.youtube.com", "youtu.be"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`v=[\w-]+`), // 精简匹配
			regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
		},
		handler: &music.YoutubeMusicHandler{},
	},
	{
		domains: []string{"music.apple.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(song|album|playlist)/[^/]+/(id)?\d+`),
		},
		handler: &music.AppleMusicHandler{},
	},
	{
		domains: []string{"www.douyin.com", "douyin.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`^https?://www\.douyin\.com/video/[\w-]+`),
		},
		handler: &video.DouYinVideoHandler{},
	},
	{
		domains: []string{"soundcloud.com", "snd.sc"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/[\w-]+(/sets)?/[\w-]+`),
		},
		handler: &music.SoundCloudHandler{},
	},
	{
		domains: []string{"y.qq.com", "c.y.qq.com", "m.y.qq.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(song|playlist|album)`),
		},
		handler: &music.QQHandler{},
	},
	{
		domains: []string{"open.spotify.com", "play.spotify.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(track|album|playlist)/[\w-]+`),
		},
		handler: &music.SpotifyHandler{},
	},
}

// 快速域名 -> matcher 索引映射表（加速匹配）
var matcherMap = make(map[string]*platformMatcher)

func init() {
	for i := range platforms {
		p := &platforms[i]
		for _, d := range p.domains {
			matcherMap[d] = p
		}
	}
}

// ParseLink ⚡ 优化，加快路径匹配，减少正则消耗
func ParseLink(text string) (string, any) {
	raw := genericURLRegex.FindString(text)
	if raw == "" {
		return "", nil
	}

	raw = cleanURLTrailingChars(raw)
	u, err := url.Parse(raw)
	if err != nil {
		return "", nil
	}

	host := strings.ToLower(u.Host)
	handler, ok := quickMatch(host, u)
	if ok {
		return raw, handler
	}
	// fallback: 正则穷举匹配
	for i := range platforms {
		for _, r := range platforms[i].patterns {
			if r.MatchString(raw) {
				return raw, platforms[i].handler
			}
		}
	}
	return "", nil
}

// 通用 URL 提取
var genericURLRegex = regexp.MustCompile(`https?://[^\s<>"'()]*[\w/#?=&-]`)

// ⚡ 优化版 Trim
func cleanURLTrailingChars(s string) string {
	s = strings.TrimSpace(s)
	runes := []rune(s)
	end := len(runes)
	for end > 0 {
		r := runes[end-1]
		if unicode.IsSpace(r) || strings.ContainsRune(".,!:;\"'()`[]{}", r) {
			end--
			continue
		}
		if r > 127 && !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			end--
			continue
		}
		break
	}
	return string(runes[:end])
}

// quickMatch 先基于域名进行快速判断
func quickMatch(host string, u *url.URL) (any, bool) {
	if p, ok := matcherMap[host]; ok {
		// 再进行一次轻量正则或路径判断
		for _, re := range p.patterns {
			if re.MatchString(u.String()) {
				return p.handler, true
			}
		}
	}
	return nil, false
}
