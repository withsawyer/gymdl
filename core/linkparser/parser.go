package linkparser

import (
	"net/url"
	"regexp"
	"strings"
	"unicode"

	"github.com/nichuanfang/gymdl/core/domain"
)

// 链接解析器

// linkTypeMatcher 处理器匹配规则
type linkTypeMatcher struct {
	patterns []*regexp.Regexp
	linkType domain.LinkType
	domains  []string // 快速判定域名
}

/* ---------------------- 变量区 ---------------------- */

// 处理器规则集
var linkTypeMatchers = []linkTypeMatcher{
	{
		/* ---------------------- 音乐平台 ---------------------- */
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
		linkType: domain.LinkNetEase,
	},
	{
		domains: []string{"youtube.com", "music.youtube.com", "youtu.be"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`v=[\w-]+`), // 精简匹配
			regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
		},
		linkType: domain.LinkYoutubeMusic,
	},
	{
		domains: []string{"music.apple.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(song|album|playlist)/[^/]+/(id)?\d+`),
		},
		linkType: domain.LinkAppleMusic,
	},
	{
		domains: []string{"soundcloud.com", "snd.sc"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/[\w-]+(/sets)?/[\w-]+`),
		},
		linkType: domain.LinkSoundcloud,
	},
	{
		domains: []string{"y.qq.com", "c.y.qq.com", "m.y.qq.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(song|playlist|album)`),
		},
		linkType: domain.LinkQQMusic,
	},
	{
		domains: []string{"open.spotify.com", "play.spotify.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(track|album|playlist)/[\w-]+`),
		},
		linkType: domain.LinkSpotify,
	},
	/* ---------------------- 视频平台(待补充) ---------------------- */
}

// 快速域名 -> matcher 索引映射表（加速匹配）
var matcherMap = make(map[string]*linkTypeMatcher)

// 通用 URL 提取
var genericURLRegex = regexp.MustCompile(`https?://[^\s<>"'()]*[\w/#?=&-]`)

/* ---------------------- 解析器初始化 ---------------------- */

// 初始化
func init() {
	for i := range linkTypeMatchers {
		l := &linkTypeMatchers[i]
		for _, d := range l.domains {
			matcherMap[d] = l
		}
	}
}

/* ---------------------- 核心方法 ---------------------- */

// ⚡ ParseLink：解析链接
func ParseLink(text string) (string, domain.LinkType) {
	raw := genericURLRegex.FindString(text)
	if raw == "" {
		return "", domain.LinkUnknown
	}
	// 链接清洗
	raw = cleanURLTrailingChars(raw)
	u, err := url.Parse(raw)
	if err != nil {
		return "", domain.LinkUnknown
	}

	host := strings.ToLower(u.Host)
	linkType, ok := quickMatch(host, u)
	if ok {
		return raw, linkType
	}

	// fallback: 正则穷举匹配
	for i := range linkTypeMatchers {
		for _, r := range linkTypeMatchers[i].patterns {
			if r.MatchString(raw) {
				return raw, linkTypeMatchers[i].linkType
			}
		}
	}
	return "", domain.LinkUnknown
}

/* ---------------------- 辅助方法 ---------------------- */

// ⚡Trim
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
func quickMatch(host string, u *url.URL) (domain.LinkType, bool) {
	if l, ok := matcherMap[host]; ok {
		// 再进行一次轻量正则或路径判断
		for _, re := range l.patterns {
			if re.MatchString(u.String()) {
				return l.linkType, true
			}
		}
	}
	return domain.LinkUnknown, false
}
