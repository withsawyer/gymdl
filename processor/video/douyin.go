package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/nichuanfang/gymdl/core"
	"github.com/nichuanfang/gymdl/processor"
	"github.com/nichuanfang/gymdl/utils"
	"github.com/playwright-community/playwright-go"
	"github.com/withsawyer/gopher-tools/datetime"
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
)

// Platform 表示平台类型
type Platform string

// DouYinProcessor 抖音视频处理器，实现视频下载功能
type DouYinProcessor struct {
	cfg       *config.Config
	tempDir   string
	videos    []*VideoInfo
	videoInfo *VideoInfo
	reporter  ProgressReporter // 进度报告器作为结构体字段
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
func (p *DouYinProcessor) Download(link string, reporter ProgressReporter) error {
	// 保存reporter到结构体字段
	p.reporter = reporter
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
	//创建通道用于接收API响应数据
	apiDataChan := make(chan map[string]interface{}, 10)
	// 拦截网络请求，尝试直接获取API响应数据
	page.On("response", func(response playwright.Response) {
		responseURL := response.URL()
		// 检查是否为抖音视频数据相关的API请求
		if strings.Contains(responseURL, "/aweme/v1/play/") ||
			strings.Contains(responseURL, "/aweme/v1/aweme/detail/") ||
			strings.Contains(responseURL, "video/play/") {
			utils.DebugWithFormat("捕获到视频相关API请求: %s", responseURL)
			// 尝试获取JSON响应
			if response.Status() == 200 {
				// 修复JSON解析方法调用
				var jsonData interface{}
				err = response.JSON(&jsonData)
				if err == nil {
					if dataMap, ok := jsonData.(map[string]interface{}); ok {
						apiDataChan <- dataMap
					}
				}
			}
		}
	})

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
		err = playwright.Install()
		if err != nil {
			return nil, nil, nil, fmt.Errorf("启动 Playwright 失败: %v", err)
		}
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
	utils.DebugWithFormat("随机选择用户代理: %s", selectedUserAgent)

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
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36 Edg/142.0.0.0",
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

	// 访问 URL - 等待网络空闲状态以确保页面完全加载
	if _, err := page.Goto(link, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle, // 等待网络空闲
		Timeout:   playwright.Float(60000),
	}); err != nil {
		return "", fmt.Errorf("访问页面失败: %v", err)
	}

	// 确保页面完全加载 - 等待所有资源加载完成
	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle, // 再次确认网络空闲
	}); err != nil {
		utils.InfoWithFormat("等待页面网络空闲失败: %v", err)
	}

	// 等待特定元素出现，确保关键内容已加载
	waitForElements(page)

	// 尝试滚动页面，确保动态加载的内容也被加载
	if _, err := page.Evaluate(`window.scrollTo(0, document.body.scrollHeight);`); err != nil {
		utils.InfoWithFormat("页面滚动失败: %v", err)
	}

	// 滚动后再次等待网络空闲
	if err := page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateNetworkidle,
	}); err != nil {
		utils.InfoWithFormat("滚动后等待网络空闲失败: %v", err)
	}

	// 使用智能等待替代硬编码延时
	if err := p.waitForVideoContent(page); err != nil {
		utils.InfoWithFormat("等待视频内容超时: %v", err)
	}

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

	// 查找包含视频数据的script标签
	scriptRegex := regexp.MustCompile(scriptPattern)
	scriptMatches := scriptRegex.FindAllStringSubmatch(html, -1)

	for _, scriptMatch := range scriptMatches {
		scriptContent := scriptMatch[1]
		// 检查是否包含关键数据标记
		if !strings.Contains(scriptContent, "aweme_id") || !strings.Contains(scriptContent, "video") {
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

// waitForVideoContent 智能等待视频内容加载完成
func (p *DouYinProcessor) waitForVideoContent(page playwright.Page) error {
	// 使用轮询检查页面是否包含视频关键数据
	deadline := time.Now().Add(30 * time.Second) // 最多等待30秒
	for time.Now().Before(deadline) {
		html, err := page.Content()
		if err == nil && (strings.Contains(html, "aweme_id") || strings.Contains(html, "video")) {
			utils.DebugWithFormat("检测到视频内容已加载")
			return nil
		}
		time.Sleep(500 * time.Millisecond) // 每500ms检查一次
	}
	return errors.New("等待视频内容超时")
}

// waitForElements 等待关键元素出现
func waitForElements(page playwright.Page) {
	// 尝试等待几个关键元素出现，但不阻塞主流程
	go func() {
		// 等待视频容器元素
		if _, err := page.WaitForSelector("title", playwright.PageWaitForSelectorOptions{
			Timeout: playwright.Float(60000),
		}); err != nil {
			utils.ErrorWithFormat("等待视频元素超时: %v", err)
		}
	}()
}

// findDataInJson 在数据结构中查找视频URL
func (p *DouYinProcessor) findDataInJson(data map[string]interface{}) *htmlJsonData {
	var hjd = &htmlJsonData{}

	// 安全提取字段，避免 panic
	getString := func(m map[string]interface{}, key string) string {
		if v, ok := m[key].(string); ok {
			return v
		}
		return ""
	}

	getMap := func(m map[string]interface{}, key string) map[string]interface{} {
		if key == "" {
			return m
		}
		if v, ok := m[key].(map[string]interface{}); ok {
			return v
		}
		return nil
	}

	getList := func(m map[string]interface{}, key string) []interface{} {
		if v, ok := m[key].([]interface{}); ok {
			return v
		}
		return nil
	}

	// 从 data 开始逐层解析
	if loaderData := getMap(data, "loaderData"); loaderData != nil {
		if videoPage := getMap(loaderData, "video_(id)/page"); videoPage != nil {
			if videoInfoRes := getMap(videoPage, "videoInfoRes"); videoInfoRes != nil {
				if itemList := getList(videoInfoRes, "item_list"); len(itemList) > 0 {
					if item := getMap(itemList[0].(map[string]interface{}), ""); item != nil {
						hjd.Desc = getString(item, "desc")

						if author := getMap(item, "author"); author != nil {
							hjd.Author = getString(author, "nickname")
						}

						if createTime, ok := item["create_time"].(int64); ok {
							hjd.Time = datetime.FormatTimeToStr(time.Unix(createTime, 0), "yyyy-mm-dd hh:mm:ss")
						}

						if video := getMap(item, "video"); video != nil {
							if playAddr := getMap(video, "play_addr"); playAddr != nil {
								if urls := getList(playAddr, "url_list"); len(urls) > 0 {
									for _, u := range urls {
										vUrl := u.(string)
										if !p.isURLAccessible(vUrl) {
											utils.InfoWithFormat("[findDataInJson] 视频URL不可访问: %s", vUrl)
											continue
										}
										hjd.VideoUrl = vUrl
										if hjd.VideoUrl != "" {
											hjd.VideoUrl = strings.Replace(hjd.VideoUrl, "playwm", "play", 1)
										}
									}

								}
							}

							if cover := getMap(video, "cover"); cover != nil {
								if urls := getList(cover, "url_list"); len(urls) > 0 {
									for _, u := range urls {
										imgUrl := u.(string)
										if !p.isURLAccessible(imgUrl) {
											utils.InfoWithFormat("[findDataInJson] 封面URL不可访问: %s", imgUrl)
											continue
										}
										hjd.CoverUrl = imgUrl
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return hjd
}

// resolveFinalURL 解析URL的最终重定向地址
func (p *DouYinProcessor) resolveFinalURL(originalURL string) (string, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 限制重定向次数，避免无限循环
			if len(via) >= 10 {
				return errors.New("too many redirects")
			}
			// 保留原始请求头
			for key, values := range via[0].Header {
				if req.Header.Get(key) == "" {
					for _, value := range values {
						req.Header.Add(key, value)
					}
				}
			}
			return nil
		},
	}

	req, err := http.NewRequest("GET", originalURL, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置请求头，模拟真实浏览器
	req.Header.Set("User-Agent", p.getRandomUserAgent())
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Connection", "keep-alive")

	// 对抖音域名添加特殊处理
	if strings.Contains(originalURL, "douyin.com") || strings.Contains(originalURL, "snssdk.com") {
		req.Header.Set("Referer", "https://www.iesdouyin.com/")
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()

	// 记录重定向信息
	if finalURL != originalURL {
		utils.DebugWithFormat("[redirect] URL重定向: %s -> %s", originalURL, finalURL)
	}

	// 检查最终状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("最终URL返回错误状态码: %d", resp.StatusCode)
	}

	return finalURL, nil
}

// isURLAccessible 检查URL是否可正常访问（增强版，支持重定向解析）
func (p *DouYinProcessor) isURLAccessible(url string) bool {
	// 首先尝试解析最终URL
	finalURL, err := p.resolveFinalURL(url)
	if err != nil {
		utils.DebugWithFormat("[accessibility] URL解析失败: %v", err)
		return false
	}

	// 使用HEAD请求检查最终URL的可用性
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("HEAD", finalURL, nil)
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
	accessible := resp.StatusCode >= 200 && resp.StatusCode < 300
	if accessible {
		utils.DebugWithFormat("[accessibility] URL可访问: %s (状态码: %d)", finalURL, resp.StatusCode)
	} else {
		utils.DebugWithFormat("[accessibility] URL不可访问: %s (状态码: %d)", finalURL, resp.StatusCode)
	}

	return accessible
}

// downloadResource 下载资源到指定路径
func (p *DouYinProcessor) downloadVideo() error {
	// 下载视频逻辑
	for _, videoInfo := range p.videos {
		// 首先下载视频文件
		fp := p.tempDir
		fn := videoInfo.Desc + ".mp4"
		downloadSize, err := p._downloadResource(videoInfo.DownloadUrl, fp, fn)
		if err != nil {
			utils.ErrorWithFormat("[download] 视频下载失败: %v", err)
			return err
		}
		videoInfo.Size = downloadSize
		utils.InfoWithFormat("[download] 视频下载完成: %s", filepath.Join(fp, fn))

		if videoInfo.CoverUrl != "" {
			// 下载封面图片，使用正确的CoverUrl
			fn = videoInfo.Desc + ".png"
			_, err = p._downloadResource(videoInfo.CoverUrl, fp, fn)
			if err != nil {
				utils.ErrorWithFormat("[download] 封面下载失败: %v", err)
				// 封面下载失败不影响整体流程，继续执行
				continue
			}
			utils.InfoWithFormat("[download] 下载完成: %s", filepath.Join(fp, fn))
		}
	}
	return nil
}

func (p *DouYinProcessor) _downloadResource(url, savePath, filename string) (string, error) {
	// 确保保存目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		return "", fmt.Errorf("创建保存目录失败: %w", err)
	}

	// 正确设置下载选项
	downloader, err := utils.DownloadFile(url, &utils.DownloadOptions{
		SavePath:   savePath,
		FileName:   filename,
		Timeout:    1200 * time.Second, // 正确设置为time.Duration类型
		IgnoreSSL:  true,
		MaxRetries: 3,                // 增加重试次数
		ChunkSize:  10 * 1024 * 1024, // 正确设置为10MB字节数
	})

	// 启动下载
	if err = downloader.Start(); err != nil {
		utils.ErrorWithFormat("[downloader] 启动下载失败: %v", err)
		return "", err
	}

	// 等待下载完成
	for {
		progress := downloader.GetProgress()
		if progress.Status == utils.StatusCompleted || progress.Status == utils.StatusFailed {
			if progress.Status == utils.StatusCompleted {
				// 返回格式化的文件大小
				return progress.FormattedSize, nil
			} else {
				// 返回具体的错误信息
				return "", fmt.Errorf("[downloader] 下载失败: %s, 错误: %s", url, progress.ErrorMessage)
			}
		}
		// 修正进度显示格式，使用Progress而不是FormattedSpeed
		progressText := fmt.Sprintf("正在下载:\n文件名%s：\n进度: %.2f%%\n速度: %s\n已下载: %s/%s",
			filename,
			progress.Progress,
			progress.FormattedSpeed,
			progress.FormattedDownloaded, progress.FormattedSize)
		if p.reporter != nil {
			p.reporter.ReportProgress(progressText)
		}
		time.Sleep(1 * time.Second) // 减少轮询间隔
	}
}

func (p *DouYinProcessor) Tidy() error {
	files, err := os.ReadDir(p.tempDir)
	if err != nil {
		return fmt.Errorf("读取临时目录失败: %w", err)
	}
	if len(files) == 0 {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 未找到待整理的资源文件")
		return errors.New("未找到待整理的资源文件")
	}

	switch p.cfg.Tidy.Mode {
	case 1:
		return p.tidyToLocal(files)
	case 2:
		return p.tidyToWebDAV(files, core.GlobalWebDAV)
	default:
		return fmt.Errorf("未知整理模式: %d", p.cfg.Tidy.Mode)
	}
}

// 整理到本地
func (p *DouYinProcessor) tidyToLocal(files []os.DirEntry) error {
	dstDir := p.cfg.Tidy.DistDir
	if dstDir == "" {
		_ = processor.RemoveTempDir(p.tempDir)
		return errors.New("未配置输出目录")
	}
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		_ = processor.RemoveTempDir(p.tempDir)
		return fmt.Errorf("创建输出目录失败: %w", err)
	}

	for _, f := range files {
		src := filepath.Join(p.tempDir, f.Name())
		mvDir := filepath.Join(dstDir, "douyin")
		if err := os.MkdirAll(mvDir, 0755); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
		dst := filepath.Join(mvDir, utils.SanitizeFileName(f.Name()))
		if err := os.Rename(src, dst); err != nil {
			utils.WarnWithFormat("[DouYinVideo] ⚠️ 移动失败 %s → %s: %v", src, dst, err)
			continue
		}
		utils.InfoWithFormat("[DouYinVideo] 📦 已整理: %s", dst)
	}
	//清除临时目录
	err := processor.RemoveTempDir(p.tempDir)
	if err != nil {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 删除临时目录失败: %s (%v)", p.tempDir, err)
		return err
	}
	utils.DebugWithFormat("[DouYinVideo] 🧹 已删除临时目录: %s", p.tempDir)
	return nil
}

// 整理到webdav
func (p *DouYinProcessor) tidyToWebDAV(files []os.DirEntry, webdav *core.WebDAV) error {
	if webdav == nil {
		_ = processor.RemoveTempDir(p.tempDir)
		return errors.New("WebDAV 未初始化")
	}

	for _, f := range files {
		filePath := filepath.Join(p.tempDir, "douyin", f.Name())
		if err := webdav.Upload(filePath); err != nil {
			utils.WarnWithFormat("[DouYinVideo] ☁️ 上传失败 %s: %v", f.Name(), err)
			continue
		}
		utils.InfoWithFormat("[DouYinVideo] ☁️ 已上传: %s", f.Name())
	}
	//清除临时目录
	err := processor.RemoveTempDir(p.tempDir)
	if err != nil {
		utils.WarnWithFormat("[DouYinVideo] ⚠️ 删除临时目录失败: %s (%v)", p.tempDir, err)
		return err
	}
	utils.DebugWithFormat("[DouYinVideo] 🧹 已删除临时目录: %s", p.tempDir)
	return nil
}
