package handler

import (
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"unicode"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
)

// 音乐信息
type SongInfo struct {
	SongName    string //音乐名称
	SongArtists string
	SongAlbum   string
	FileExt     string //格式
	MusicSize   int    //音乐大小
	Bitrate     string //码率
	Duration    int    //时长
	PicUrl      string //封面图url
	Tidy        string //入库方式(默认/webdav)
}

// MusicHandler 音乐处理接口
type MusicHandler interface {
	//平台名称
	Platform() string
	//下载音乐
	DownloadMusic(url string, cfg *config.Config) (*SongInfo, error)
	//构建下载命令
	DownloadCommand(cfg *config.Config, url string) *exec.Cmd
	//音乐整理之前的处理(如嵌入元数据,刮削等)
	BeforeTidy(cfg *config.Config, songInfo *SongInfo) error
	//是否需要移除DRM
	NeedRemoveDRM(cfg *config.Config) bool
	//移除DRM
	DRMRemove(cfg *config.Config, songInfo *SongInfo) error
	//音乐整理
	TidyMusic(cfg *config.Config, webdav *core.WebDAV, songInfo *SongInfo) error
	//加密后缀
	EncryptedExts() []string
	//非加密后缀
	DecryptedExts() []string
}

// platformMatcher 平台匹配规则
type platformMatcher struct {
	patterns []*regexp.Regexp
	handler  MusicHandler
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
			//其他平台
			regexp.MustCompile(`^https?://163cn\.link/\w+`),
		},
		handler: &NCMHandler{},
	},
	{
		domains: []string{"youtube.com", "music.youtube.com", "youtu.be"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`v=[\w-]+`), // 精简匹配
			regexp.MustCompile(`^https?://youtu\.be/[\w-]+`),
		},
		handler: &YoutubeMusicHandler{},
	},
	{
		domains: []string{"music.apple.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(song|album|playlist)/[^/]+/(id)?\d+`),
		},
		handler: &AppleMusicHandler{},
	},
	{
		domains: []string{"soundcloud.com", "snd.sc"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/[\w-]+(/sets)?/[\w-]+`),
		},
		handler: &SoundCloudHandler{},
	},
	{
		domains: []string{"y.qq.com", "c.y.qq.com", "m.y.qq.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(song|playlist|album)`),
		},
		handler: &QQHandler{},
	},
	{
		domains: []string{"open.spotify.com", "play.spotify.com"},
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`/(track|album|playlist)/[\w-]+`),
		},
		handler: &SpotifyHandler{},
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

// ⚡ 优化 ParseLink：加快路径匹配，减少正则消耗
func ParseLink(text string) (string, MusicHandler) {
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

// quickMatch 先基于域名进行快速判断
func quickMatch(host string, u *url.URL) (MusicHandler, bool) {
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
