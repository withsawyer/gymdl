package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
	"github.com/playwright-community/playwright-go"
	"github.com/withsawyer/gopher-tools/datetime"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	// videoURLPattern 视频URL匹配模式
	videoURLPattern = `/video/(\d+)`
	// videoIDPattern 视频ID匹配模式
	videoIDPattern = `video_id=([a-zA-Z0-9]+)`
	// jsonPattern JSON数据匹配模式
	jsonPattern = `({.*?"errors":\s*null\s*})`
	// scriptPattern 脚本标签匹配模式
	scriptPattern = `<script[^>]*>(.*?)</script>`
	// titlePattern 标题匹配模式
	titlePattern = `<title>(.*?)</title>`
)

// Platform 表示平台类型
type Platform string

// DouYinProcessor 抖音视频处理器，实现视频下载功能
type DouYinProcessor struct {
	cfg       *config.Config
	tempDir   string
	videos    []*VideoInfo
	videoInfo *VideoInfo
}

// Init 初始化抖音处理器
func (p *DouYinProcessor) Init(cfg *config.Config) {
	p.cfg = cfg
	p.videos = make([]*VideoInfo, 0)
	p.tempDir = processor.BuildOutputDir(DouyinTempDir)
	p.videoInfo = &VideoInfo{}
}

// Name 返回处理器名称
func (p *DouYinProcessor) Name() processor.LinkType {
	return processor.LinkDouyin
}

// Videos 返回已下载的视频信息列表
func (p *DouYinProcessor) Videos() []*VideoInfo {
	return p.videos
}

// Download 下载抖音视频
func (p *DouYinProcessor) Download(link string) error {
	err := p.method1(link)
	if err != nil {
		return fmt.Errorf("下载抖音视频失败: %v", err)
	}

	return nil
}

func (p *DouYinProcessor) method1(link string) error {
	// 初始化 Playwright 和浏览器
	ctx, page, pw, err := p.initPlaywrightAndBrowser()
	if err != nil {
		return err
	}
	defer func() {
		page.Close()
		ctx.Close()
		pw.Stop()
	}()

	// 加载 cookies
	if err = p.loadCookies(ctx); err != nil {
		utils.InfoWithFormat("加载 cookies 失败: %v", err)
	}

	// 提取视频ID
	videoID, err := p._extractVideoID(page, link)
	if err != nil {
		return err
	}
	utils.InfoWithFormat("提取视频ID成功: %s", videoID)
	// 提取视频内容和URL
	html, err := page.Content()
	if err != nil {
		return fmt.Errorf("获取页面内容失败: %v", err)
	}

	err = p.extractDataFromHTML(html)
	// 保存视频信息，当前只获取一个视频，所以直接保存
	p.videos = append(p.videos, p.videoInfo)
	if err != nil {
		return fmt.Errorf("提取视频 URL 失败: %v", err)
	}
	// 下载视频
	if err = p.downloadVideo(); err != nil {
		return fmt.Errorf("下载视频失败: %v", err)
	}
	return nil
}

// initPlaywrightAndBrowser 初始化 Playwright 和浏览器
func (p *DouYinProcessor) initPlaywrightAndBrowser() (playwright.BrowserContext, playwright.Page, *playwright.Playwright, error) {
	pw, err := playwright.Run()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("启动 Playwright 失败: %v", err)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		pw.Stop()
		return nil, nil, nil, fmt.Errorf("启动浏览器失败: %v", err)
	}

	// 随机选择用户代理
	selectedUserAgent := p.getRandomUserAgent()

	// 创建浏览器上下文
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
		browser.Close()
		pw.Stop()
		return nil, nil, nil, fmt.Errorf("创建上下文失败: %v", err)
	}

	page, err := ctx.NewPage()
	if err != nil {
		ctx.Close()
		browser.Close()
		pw.Stop()
		return nil, nil, nil, fmt.Errorf("创建页面失败: %v", err)
	}

	return ctx, page, pw, nil
}

// getRandomUserAgent 获取随机用户代理
func (p *DouYinProcessor) getRandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (iPhone; CPU iPhone OS 16_0 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.0 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36 QuarkPC/4.6.0.558",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; SM-G975F) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
	}
	// 随机选择一条 userAgent
	rand.New(rand.NewSource(time.Now().Unix()))
	return userAgents[rand.Intn(len(userAgents))]
}

// loadCookies 加载 cookies
func (p *DouYinProcessor) loadCookies(ctx playwright.BrowserContext) error {
	cookies := p.parseDouYinCookiesFile()
	if cookies != nil && len(cookies) > 0 {
		if err := ctx.AddCookies(cookies); err != nil {
			return err
		}
		utils.InfoWithFormat("成功加载 %d 个 cookies", len(cookies))
	}
	return nil
}

// _extractVideoID 提取视频ID
func (p *DouYinProcessor) _extractVideoID(page playwright.Page, link string) (string, error) {
	videoID := ""

	// 监听网络请求中的 video_id
	page.On("request", func(request playwright.Request) {
		requestURL := request.URL()
		if strings.Contains(requestURL, "video_id=") {
			m := regexp.MustCompile(videoIDPattern).FindStringSubmatch(requestURL)
			if len(m) > 1 {
				videoID = m[1]
				utils.DebugWithFormat("网络请求中捕获到 video_id: %s", videoID)
			}
		}
	})

	// 访问 URL - 修改为等待网络空闲状态
	if _, err := page.Goto(link, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle, // 改为等待网络空闲
		Timeout:   playwright.Float(30000),
	}); err != nil {
		return "", fmt.Errorf("访问页面失败: %v", err)
	}

	// 确保页面完全加载 - 增加额外等待和滚动操作
	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	}); err != nil {
		utils.InfoWithFormat("等待页面加载完成失败: %v", err)
	}

	// 尝试滚动页面，确保动态加载的内容也被加载
	if _, err := page.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`); err != nil {
		utils.InfoWithFormat("页面滚动失败: %v", err)
	}

	// 再等待一段时间确保内容完全加载
	//time.Sleep(2 * time.Second)

	// 从当前URL提取视频ID
	currentURL := page.URL()
	if m := regexp.MustCompile(videoURLPattern).FindStringSubmatch(currentURL); len(m) > 1 {
		videoID = m[1]
		utils.DebugWithFormat("从当前 URL 提取到 video_id: %s", videoID)
	}

	// 从原始URL提取视频ID作为备选
	if videoID == "" {
		log.Println("未捕获到 video_id，尝试从 URL 直接提取")
		if m := regexp.MustCompile(videoURLPattern).FindStringSubmatch(link); len(m) > 1 {
			videoID = m[1]
			utils.DebugWithFormat("从原始 URL 提取到 video_id: %s", videoID)
		}
	}

	// 如果仍然没有获取到videoID，尝试从页面内容中搜索aweme_id
	if videoID == "" {
		utils.DebugWithFormat("尝试从页面内容中搜索aweme_id")
		html, err := page.Content()
		if err == nil {
			// 尝试直接匹配aweme_id
			awemeIDRegex := regexp.MustCompile(`"aweme_id"\s*:\s*"([^"]+)"`)
			if m := awemeIDRegex.FindStringSubmatch(html); len(m) > 1 {
				videoID = m[1]
				utils.DebugWithFormat("从页面内容中提取到 aweme_id: %s", videoID)
			}
		}
	}

	if videoID == "" {
		return "", errors.New("未能捕获到视频数据")
	}

	return videoID, nil
}

// parseDouYinCookiesFile 解析抖音 cookies 文件
func (p *DouYinProcessor) parseDouYinCookiesFile() []playwright.OptionalCookie {
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
			// 直接在这里处理cookies，避免值传递问题
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
	}
	return playwrightCookies
}

// extractDataFromHTML 从HTML中提取视频URL
func (p *DouYinProcessor) extractDataFromHTML(html string) error {
	utils.DebugWithFormat("[extract] HTML长度: %d 字符", len(html))

	// 提取视频标题
	titleMatches := regexp.MustCompile(titlePattern).FindAllStringSubmatch(html, -1)
	var title string
	for _, titleMatch := range titleMatches {
		title = titleMatch[1]
		if !strings.Contains(title, "-") {
			continue
		}
		if title != "" {
			sTitle := strings.Split(title, "-")
			p.videoInfo.Title = sTitle[0]
		}
	}
	utils.InfoWithFormat("[extract] 提取到视频标题: %s", title)

	// 查找包含视频数据的script标签
	scriptRegex := regexp.MustCompile(scriptPattern)
	scriptMatches := scriptRegex.FindAllStringSubmatch(html, -1)

	for _, scriptMatch := range scriptMatches {
		scriptContent := scriptMatch[1]
		// 检查是否包含关键数据标记
		if !strings.Contains(scriptContent, "aweme_id") || !strings.Contains(scriptContent, "status_code") {
			continue
		}
		// 尝试提取JSON部分
		jsonMatches := regexp.MustCompile(jsonPattern).FindAllStringSubmatch(scriptContent, -1)
		for _, jsonMatch := range jsonMatches {
			jsonStr := jsonMatch[1]
			// 清理JSON，确保匹配完整的JSON结构
			cleanJSON, err := p.cleanJSONString(jsonStr)
			if err != nil || cleanJSON == "" {
				continue
			}
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(cleanJSON), &data); err != nil {
				continue
			}
			// 递归查找视频URL
			hjd := p.findDataInJson(data)
			if hjd.VideoUrl != "" {
				p.videoInfo.CoverUrl = hjd.CoverUrl
				p.videoInfo.Time = hjd.Time
				p.videoInfo.Desc = hjd.Desc
				p.videoInfo.Author = hjd.Author
				p.videoInfo.Ratio = hjd.Ratio
				p.videoInfo.DownloadUrl = hjd.VideoUrl
				utils.InfoWithFormat("[extract] 提取到视频信息: %v", hjd)
				return nil
			}
		}
	}

	return errors.New("未能提取到视频URL")
}

// cleanJSONString 清理JSON字符串，确保其完整性
func (p *DouYinProcessor) cleanJSONString(jsonStr string) (string, error) {
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
		return jsonStr[:jsonEnd], nil
	}
	return "", errors.New("无法找到完整的JSON结构")
}

type htmlJsonData struct {
	VideoUrl string
	CoverUrl string
	Time     string
	Desc     string
	Author   string
	Ratio    string
}

// findDataInJson 在数据结构中查找视频URL
func (p *DouYinProcessor) findDataInJson(data map[string]interface{}) *htmlJsonData {
	var hjd = &htmlJsonData{}
	var findURL func(obj interface{}) *htmlJsonData
	findURL = func(obj interface{}) *htmlJsonData {
		switch v := obj.(type) {
		case map[string]interface{}:
			for key, value := range v {
				// 专门查找video字段
				if key == "video" {
					if videoData, ok := value.(map[string]interface{}); ok {
						utils.DebugWithFormat("[extract] 找到video字段: %v", getMapKeys(videoData))

						// 优先查找play_url字段（无水印）
						videoURL := p.extractURLFromField(videoData, "play_url", false)
						if videoURL == "" {
							// 兜底：如果没有play_url，再查找play_addr字段（有水印）
							videoURL = p.extractURLFromField(videoData, "play_addr", true)
						}
						if videoURL != "" {
							hjd.VideoUrl = videoURL
							params, err := p._extractURLParams(videoURL)
							if err == nil {
								if ratio := params["ratio"]; ratio != "" {
									hjd.Ratio = ratio
								}
							}
						}
						coverURL := p.extractURLFromField(videoData, "cover", false)
						if coverURL != "" {
							hjd.CoverUrl = coverURL
						}
					}
				} else if key == "desc" {
					hjd.Desc = value.(string)
				} else if key == "author" {
					if authorData, ok := value.(map[string]interface{}); ok {
						if nickname := authorData["nickname"].(string); nickname != "" {
							hjd.Author = nickname
						}
					}
				} else if key == "create_time" {
					ctTime := datetime.FormatTimeToStr(time.Unix(value.(int64), 0), "yyyy-mm-dd hh:mm:ss")
					hjd.Time = ctTime
				} else {
					// 递归搜索其他字段
					if result := findURL(value); result != nil {
						return result
					}
				}
				return hjd
			}
		case []interface{}:
			for _, item := range v {
				if result := findURL(item); result != nil {
					return result
				}
			}
		}
		return nil
	}

	return findURL(data)
}

// extractURLFromField 从指定字段提取URL
func (p *DouYinProcessor) extractURLFromField(data map[string]interface{}, fieldName string, isVideoWatermarked bool) string {
	field, ok := data[fieldName]
	if !ok {
		return ""
	}

	fieldMap, ok := field.(map[string]interface{})
	if !ok {
		return ""
	}

	urlList, ok := fieldMap["url_list"].([]interface{})
	if !ok || len(urlList) == 0 {
		return ""
	}

	// 遍历所有URL，跳过无效地址
	for _, item := range urlList {
		if url, ok := item.(string); ok && strings.HasPrefix(url, "http") {
			utils.DebugWithFormat("[extract] 从%s.url_list找到: %s", fieldName, url)

			// 检查URL是否能正常访问
			if !p.isURLAccessible(url) {
				utils.DebugWithFormat("[extract] URL不可访问，跳过: %s", url)
				continue
			}

			// 如果是有水印的视频，尝试去除水印
			if isVideoWatermarked && strings.Contains(url, "playwm") {
				// 替换playwm为play，转换为无水印视频
				url = strings.Replace(url, "playwm", "play", 1)
				// 检查转换后的URL是否可访问
				if !p.isURLAccessible(url) {
					utils.DebugWithFormat("[extract] 转换后的URL不可访问，使用原URL: %s", url)
					// 恢复原URL
					url = strings.Replace(url, "play", "playwm", 1)
				}
			}
			return url
		}
		utils.DebugWithFormat("[extract] 跳过无效URL: %v", item)
	}

	return ""
}

// isURLAccessible 检查URL是否可正常访问
func (p *DouYinProcessor) isURLAccessible(url string) bool {
	// 使用HEAD请求快速检查URL可用性，避免下载整个文件
	client := &http.Client{
		Timeout: 5 * time.Second, // 设置超时时间
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false
	}

	// 设置User-Agent，避免被拦截
	req.Header.Set("User-Agent", p.getRandomUserAgent())

	resp, err := client.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	// 检查响应状态码，200-299表示成功
	return resp.StatusCode >= 200 && resp.StatusCode < 300
}

// isVideoURL 判断是否为有效的视频URL
func (p *DouYinProcessor) isVideoURL(url string) bool {
	videoExtensions := []string{".mp4", ".m3u8", ".ts", "douyinvod.com", "snssdk.com"}
	for _, ext := range videoExtensions {
		if strings.Contains(strings.ToLower(url), ext) {
			return true
		}
	}
	return false
}

// getMapKeys 获取map的所有键
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// downloadVideo 下载视频到指定路径
func (p *DouYinProcessor) downloadVideo() error {
	// 下载视频逻辑
	for _, videoInfo := range p.videos {

	}
	return nil
}

// _extractURLParams 从 URL 中提取查询参数，并返回一个键值对映射。
func (p *DouYinProcessor) _extractURLParams(rawURL string) (map[string]string, error) {
	// 解析 URL
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}
	// 提取查询参数
	queryParams := parsedURL.Query()
	params := make(map[string]string)
	// 将查询参数转换为 map
	for key, values := range queryParams {
		// 如果参数有多个值，只取第一个值
		if len(values) > 0 {
			params[key] = values[0]
		}
	}
	return params, nil
}
