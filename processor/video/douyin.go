package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
	"github.com/playwright-community/playwright-go"
	"log"
	"math/rand"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Platform string

type DouYinProcessor struct {
	cfg     *config.Config
	tempDir string
	videos  []*VideoInfo
}

func (p *DouYinProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.videos = make([]*VideoInfo, 0)
	p.tempDir = processor.BuildOutputDir(DouyinTempDir)
}

func (p *DouYinProcessor) Name() processor.LinkType {
	return processor.LinkDouyin
}

func (p *DouYinProcessor) Videos() []*VideoInfo {
	return p.videos
}

func (p *DouYinProcessor) Download(link string) error {
	pw, err := p.runPlaywright()
	if err != nil {
		return errors.New(fmt.Sprintf("启动 Playwright 失败，无法下载抖音视频：%v", err))
	}
	defer pw.Stop()
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("启动浏览器失败: %v", err)
	}
	defer browser.Close()

	// 定义 userAgents 列表
	userAgents := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 QuarkPC/4.6.0.558",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	}
	// 随机选择一条 userAgent
	rand.New(rand.NewSource(time.Now().Unix()))
	selectedUserAgent := userAgents[rand.Intn(len(userAgents))]

	contextOptions := playwright.BrowserNewContextOptions{
		UserAgent:         playwright.String(selectedUserAgent),
		Viewport:          &playwright.Size{Width: 375, Height: 667},
		DeviceScaleFactor: playwright.Float(2),
		Locale:            playwright.String("zh-CN"),
		TimezoneId:        playwright.String("Asia/Shanghai"),
		IsMobile:          playwright.Bool(true),
		HasTouch:          playwright.Bool(true),
		ColorScheme:       (*playwright.ColorScheme)(playwright.String("light")),
		ExtraHttpHeaders: map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,video/mp4,*/*;q=0.8",
			"Accept-Language":           "zh-CN,zh;q=0.9,en;q=0.8",
			"Connection":                "keep-alive",
			"Upgrade-Insecure-Requests": "1",
		},
	}

	ctx, err := browser.NewContext(contextOptions)
	if err != nil {
		return fmt.Errorf("创建上下文失败: %v", err)
	}
	defer ctx.Close()

	page, err := ctx.NewPage()
	if err != nil {
		return fmt.Errorf("创建页面失败: %v", err)
	}
	defer page.Close()

	// 加载 cookies（如果存在）
	cookies := p.parseDouYinCookiesFile()
	if cookies != nil && len(cookies) > 0 {
		if err = ctx.AddCookies(cookies); err != nil {
			utils.InfoWithFormat("加载 cookies 失败: %v", err)
		} else {
			utils.InfoWithFormat("成功加载 %d 个 cookies", len(cookies))
		}
	}
	// 监听网络请求中的 video_id
	videoID := ""
	page.On("request", func(request playwright.Request) {
		requestURL := request.URL()
		if strings.Contains(requestURL, "video_id=") {
			m := regexp.MustCompile(`video_id=([a-zA-Z0-9]+)`).FindStringSubmatch(requestURL)
			if len(m) > 1 {
				videoID = m[1]
				log.Printf("网络请求中捕获到 video_id: %s", videoID)
			}
		}
	})

	// 访问 URL
	if _, err := page.Goto(link, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
		Timeout:   playwright.Float(30000),
	}); err != nil {
		return fmt.Errorf("访问页面失败: %v", err)
	}

	// 提取 video_id
	currentURL := page.URL()
	if m := regexp.MustCompile(`/video/(\d+)`).FindStringSubmatch(currentURL); len(m) > 1 {
		videoID = m[1]
		log.Printf("从当前 URL 提取到 video_id: %s", videoID)
	}

	// 等待 video_id 出现
	if videoID == "" {
		log.Println("未捕获到 video_id，尝试从 URL 直接提取")
		if m := regexp.MustCompile(`/video/(\d+)`).FindStringSubmatch(link); len(m) > 1 {
			videoID = m[1]
			log.Printf("从原始 URL 提取到 video_id: %s", videoID)
		}
	}

	if videoID == "" {
		return errors.New("未能捕获到视频数据")
	}

	// 提取视频 URL
	html, err := page.Content()
	if err != nil {
		return fmt.Errorf("获取页面内容失败: %v", err)
	}

	videoURL, err := p._extractDouYinURLFromHTML(html)
	if err != nil {
		return fmt.Errorf("提取视频 URL 失败: %v", err)
	}

	// 下载视频
	if err := p.downloadVideo(videoURL, filepath.Join(p.tempDir, videoID+".mp4")); err != nil {
		return fmt.Errorf("下载视频失败: %v", err)
	}

	return nil
}

func (p *DouYinProcessor) runPlaywright() (*playwright.Playwright, error) {
	pw, err := playwright.Run()
	return pw, err
}

func (p *DouYinProcessor) parseDouYinCookiesFile() []playwright.OptionalCookie {
	// 解析 cookies 文件逻辑
	playwrightCookies := make([]playwright.OptionalCookie, 0)
	domains := []string{
		".douyin.com",
		"douyin.com",
		"www.douyin.com",
		"v.douyin.com",
		"www.iesdouyin.com",
		"iesdouyin.com",
	}
	for _, domain := range domains {
		cookies := utils.GetCookiesByDomain(p.cfg.CookieCloud.CookieFilePath, domain)
		if len(cookies) > 0 {
			p.spliceDouYinCookies(playwrightCookies, cookies, domain)
		}
	}
	return playwrightCookies
}

func (p *DouYinProcessor) spliceDouYinCookies(playwrightCookies []playwright.OptionalCookie, cookies map[string]string, domain string) {
	for name, value := range cookies {
		playwrightCookies = append(playwrightCookies, playwright.OptionalCookie{
			Name:     name,
			Value:    value,
			Domain:   playwright.String(".douyin.com"),
			Path:     playwright.String("/"),
			HttpOnly: playwright.Bool(true),
			Secure:   playwright.Bool(true),
			SameSite: (*playwright.SameSiteAttribute)(playwright.String("Lax")),
		})
	}
}

func (p *DouYinProcessor) _extractDouYinURLFromHTML(html string) (string, error) {
	log.Printf("[extract] HTML长度: %d 字符", len(html))

	// 查找包含视频数据的script标签
	scriptRegex := regexp.MustCompile(`<script[^>]*>(.*?)</script>`)
	scriptMatches := scriptRegex.FindAllStringSubmatch(html, -1)

	for _, scriptMatch := range scriptMatches {
		scriptContent := scriptMatch[1]
		if strings.Contains(scriptContent, "aweme_id") && strings.Contains(scriptContent, "status_code") {
			// 尝试提取JSON部分
			jsonRegex := regexp.MustCompile(`({.*?"errors":\s*null\s*})`)
			jsonMatches := jsonRegex.FindAllStringSubmatch(scriptContent, -1)

			for _, jsonMatch := range jsonMatches {
				jsonStr := jsonMatch[1]
				// 清理JSON
				braceCount := 0
				jsonEnd := -1
				for i, char := range jsonStr {
					if char == '{' {
						braceCount++
					} else if char == '}' {
						braceCount--
						if braceCount == 0 {
							jsonEnd = i + 1
							break
						}
					}
				}

				if jsonEnd > 0 {
					cleanJSON := jsonStr[:jsonEnd]
					var data map[string]interface{}
					if err := json.Unmarshal([]byte(cleanJSON), &data); err != nil {
						continue
					}

					// 递归查找video字段中的视频URL
					var findVideoURL func(obj interface{}) string
					findVideoURL = func(obj interface{}) string {
						switch v := obj.(type) {
						case map[string]interface{}:
							for key, value := range v {
								// 专门查找video字段
								if key == "video" {
									if videoData, ok := value.(map[string]interface{}); ok {
										log.Printf("[extract] 找到video字段: %v", getKeys(videoData))

										// 优先查找play_url字段（无水印）
										if playURL, ok := videoData["play_url"]; ok {
											log.Printf("[extract] play_url字段内容: %v", playURL)
											log.Printf("[extract] play_url类型: %T", playURL)

											// 处理play_url字典格式
											if playURLMap, ok := playURL.(map[string]interface{}); ok {
												if urlList, ok := playURLMap["url_list"].([]interface{}); ok && len(urlList) > 0 {
													if videoURL, ok := urlList[0].(string); ok && strings.HasPrefix(videoURL, "http") {
														log.Printf("[extract] 从play_url.url_list找到无水印视频URL: %s", videoURL)
														return videoURL
													}
												}
											} else if playURLStr, ok := playURL.(string); ok && strings.HasPrefix(playURLStr, "http") {
												if containsAny(playURLStr, []string{".mp4", ".m3u8", ".ts", "douyinvod.com", "snssdk.com"}) {
													log.Printf("[extract] 找到无水印视频URL: %s", playURLStr)
													return playURLStr
												}
											}
										}

										// 兜底：如果没有play_url，再查找play_addr字段（有水印）
										if playAddr, ok := videoData["play_addr"]; ok {
											log.Printf("[extract] play_addr字段内容: %v", playAddr)
											log.Printf("[extract] play_addr类型: %T", playAddr)

											// 处理play_addr字典格式
											if playAddrMap, ok := playAddr.(map[string]interface{}); ok {
												if urlList, ok := playAddrMap["url_list"].([]interface{}); ok && len(urlList) > 0 {
													if videoURL, ok := urlList[0].(string); ok && strings.HasPrefix(videoURL, "http") {
														log.Printf("[extract] 从play_addr.url_list找到有水印视频URL: %s", videoURL)
														if strings.Contains(videoURL, "playwm") {
															// 替换playwm为play，转换为无水印视频
															videoURL = strings.Replace(videoURL, "playwm", "play", 1)
														}
														return videoURL
													}
												}
											} else if playAddrStr, ok := playAddr.(string); ok && strings.HasPrefix(playAddrStr, "http") {
												if containsAny(playAddrStr, []string{".mp4", ".m3u8", ".ts", "douyinvod.com", "snssdk.com"}) {
													log.Printf("[extract] 找到有水印视频URL: %s", playAddrStr)
													return playAddrStr
												}
											}
										}
									}
								} else {
									if result := findVideoURL(value); result != "" {
										return result
									}
								}
							}
						case []interface{}:
							for _, item := range v {
								if result := findVideoURL(item); result != "" {
									return result
								}
							}
						}
						return ""
					}

					videoURL := findVideoURL(data)
					if videoURL != "" {
						return videoURL, nil
					}
				}
			}
		}
	}

	return "", errors.New("未能从HTML中提取视频URL")
}

func getKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func containsAny(s string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(strings.ToLower(s), substr) {
			return true
		}
	}
	return false
}

func (p *DouYinProcessor) downloadVideo(url, savePath string) error {
	// 下载视频逻辑
	return nil
}
