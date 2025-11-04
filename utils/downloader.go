package utils

import (
	"bufio"
	"crypto/md5"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// 下载状态常量
type DownloadStatus int

const (
	StatusWaiting DownloadStatus = iota
	StatusDownloading
	StatusPaused
	StatusCompleted
	StatusFailed
)

// 下载错误码
type DownloadErrorCode int

const (
	ErrCodeInvalidURL DownloadErrorCode = iota + 1001
	ErrCodeNetworkError
	ErrCodeFileWriteError
	ErrCodePermissionError
	ErrCodeChecksumError
	ErrCodeTimeout
	ErrCodeInvalidConfig
	ErrCodeResumeNotSupported
)

// 下载错误信息
type DownloadError struct {
	Code    DownloadErrorCode
	Type    string
	Message string
	Err     error
}

func (e *DownloadError) Error() string {
	return fmt.Sprintf("下载错误 [代码:%d 类型:%s]: %s, 详细错误: %v", e.Code, e.Type, e.Message, e.Err)
}

// 下载进度信息
type DownloadProgress struct {
	URL                 string  // 下载URL
	TotalBytes          int64   // 总字节数
	Downloaded          int64   // 已下载字节数
	Progress            float64 // 下载进度百分比(精确到小数点后一位)
	Speed               float64 // 当前下载速度(KB/s 或 MB/s)
	Status              DownloadStatus
	ErrorMessage        string
	LastUpdate          time.Time
	FormattedSize       string // 格式化后的总大小
	FormattedDownloaded string // 格式化后的已下载大小
	FormattedSpeed      string // 格式化后的下载速度
}

// 下载选项配置
type DownloadOptions struct {
	SavePath      string                  // 保存路径
	FileName      string                  // 保存文件名
	Timeout       time.Duration           // 请求超时时间(默认30秒)
	IgnoreSSL     bool                    // 是否忽略SSL验证
	ChecksumType  string                  // 校验和类型: "md5" | "sha256" | ""(不校验)
	ExpectedHash  string                  // 期望的校验和值
	ProgressFunc  func(*DownloadProgress) // 进度回调函数
	MaxRetries    int                     // 最大重试次数(默认3次)
	RetryInterval time.Duration           // 重试间隔(默认5秒)
	ChunkSize     int                     // 分块大小(默认4MB)
}

// 下载管理器
type DownloadManager struct {
	client     *http.Client
	progress   *DownloadProgress
	options    *DownloadOptions
	stopChan   chan struct{}
	pauseChan  chan struct{}
	resumeChan chan struct{}
	mutex      sync.RWMutex
	startTime  time.Time
	lastBytes  int64
	chunkSize  int
	retries    int
}

// 创建默认下载选项
func NewDefaultDownloadOptions() *DownloadOptions {
	return &DownloadOptions{
		Timeout:       30 * time.Second,
		IgnoreSSL:     false,
		ChecksumType:  "",
		ExpectedHash:  "",
		MaxRetries:    3,
		RetryInterval: 5 * time.Second,
		ChunkSize:     4 * 1024 * 1024, // 4MB
	}
}

// 下载文件的主函数
func NewDownloader(url string, options *DownloadOptions) (*DownloadManager, error) {
	// 使用默认选项如果传入nil
	if options == nil {
		options = NewDefaultDownloadOptions()
	}

	// 验证URL
	if err := validateURL(url); err != nil {
		return nil, &DownloadError{
			Code:    ErrCodeInvalidURL,
			Type:    "URL验证错误",
			Message: "无效的下载链接",
			Err:     err,
		}
	}

	// 创建HTTP客户端
	client := createHTTPClient(options.Timeout, options.IgnoreSSL)

	// 初始化下载管理器
	manager := &DownloadManager{
		client:     client,
		options:    options,
		stopChan:   make(chan struct{}),
		pauseChan:  make(chan struct{}),
		resumeChan: make(chan struct{}, 1), // 使用缓冲通道避免阻塞
		progress: &DownloadProgress{
			URL:        url,
			Status:     StatusWaiting,
			TotalBytes: 0,
			Downloaded: 0,
			Progress:   0.0,
			Speed:      0.0,
			LastUpdate: time.Now(),
		},
		chunkSize: options.ChunkSize,
	}

	return manager, nil
}

// 开始下载
func (dm *DownloadManager) Start() error {
	dm.mutex.Lock()
	if dm.progress.Status == StatusDownloading {
		dm.mutex.Unlock()
		return errors.New("下载已经在进行中")
	}
	dm.progress.Status = StatusDownloading
	dm.startTime = time.Now()
	dm.lastBytes = 0
	dm.retries = 0
	// 初始化格式化字段
	dm.progress.FormattedSize = FormatBytes(dm.progress.TotalBytes)
	dm.progress.FormattedDownloaded = FormatBytes(dm.progress.Downloaded)
	dm.progress.FormattedSpeed = FormatSpeed(dm.progress.Speed)
	dm.mutex.Unlock()

	// 异步执行下载
	go func() {
		if err := dm.doDownload(); err != nil {
			dm.mutex.Lock()
			dm.progress.Status = StatusFailed
			dm.progress.ErrorMessage = err.Error()
			dm.mutex.Unlock()

			if dm.options.ProgressFunc != nil {
				dm.options.ProgressFunc(dm.progress)
			}
		}
	}()

	return nil
}

// 暂停下载
func (dm *DownloadManager) Pause() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	if dm.progress.Status != StatusDownloading {
		return errors.New("当前不在下载状态，无法暂停")
	}

	dm.progress.Status = StatusPaused
	select {
	case dm.pauseChan <- struct{}{}:
		return nil
	default:
		return nil // 避免通道阻塞
	}
}

// 继续下载
func (dm *DownloadManager) Resume() error {
	dm.mutex.Lock()
	if dm.progress.Status != StatusPaused {
		dm.mutex.Unlock()
		return errors.New("当前不在暂停状态，无法继续")
	}
	dm.progress.Status = StatusDownloading
	dm.mutex.Unlock()

	// 发送恢复信号
	select {
	case dm.resumeChan <- struct{}{}:
		// 重新启动下载
		go func() {
			if err := dm.doDownload(); err != nil {
				dm.mutex.Lock()
				dm.progress.Status = StatusFailed
				dm.progress.ErrorMessage = err.Error()
				dm.mutex.Unlock()

				if dm.options.ProgressFunc != nil {
					dm.options.ProgressFunc(dm.progress)
				}
			}
		}()
		return nil
	default:
		return nil // 避免通道阻塞
	}
}

// 停止下载
func (dm *DownloadManager) Stop() error {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	if dm.progress.Status != StatusDownloading && dm.progress.Status != StatusPaused {
		return errors.New("当前不在下载或暂停状态，无法停止")
	}

	dm.progress.Status = StatusFailed
	dm.progress.ErrorMessage = "下载被用户终止"

	// 关闭通道
	close(dm.stopChan)

	return nil
}

// 获取当前进度
func (dm *DownloadManager) GetProgress() *DownloadProgress {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	// 返回副本避免并发修改问题
	progress := *dm.progress
	// 确保格式化字段是最新的
	progress.FormattedSize = FormatBytes(progress.TotalBytes)
	progress.FormattedDownloaded = FormatBytes(progress.Downloaded)
	progress.FormattedSpeed = FormatSpeed(progress.Speed * 1024) // 转换回字节/秒
	return &progress
}

// 执行实际下载
func (dm *DownloadManager) doDownload() error {
	// 获取文件信息，确定保存路径和文件名
	savePath, err := dm.getSavePath()
	if err != nil {
		return &DownloadError{
			Code:    ErrCodeFileWriteError,
			Type:    "文件路径错误",
			Message: "无法确定文件保存路径",
			Err:     err,
		}
	}

	// 检查是否已有部分下载的文件
	existingSize, resumeSupported := dm.checkExistingFile(savePath)
	dm.mutex.Lock()
	dm.progress.Downloaded = existingSize
	dm.mutex.Unlock()

	// 发送初始进度更新
	if dm.options.ProgressFunc != nil {
		go dm.options.ProgressFunc(dm.GetProgress())
	}

	for {
		// 尝试下载
		err := dm.downloadWithResume(savePath, existingSize, resumeSupported)
		if err == nil {
			// 下载成功，检查完整性
			if err := dm.verifyDownload(savePath); err != nil {
				return err
			}

			dm.mutex.Lock()
			dm.progress.Status = StatusCompleted
			dm.progress.Progress = 100.0
			dm.mutex.Unlock()

			if dm.options.ProgressFunc != nil {
				dm.options.ProgressFunc(dm.GetProgress())
			}

			return nil
		}

		// 处理错误
		var downloadErr *DownloadError
		if errors.As(err, &downloadErr) {
			// 如果是停止信号，直接返回
			if dm.progress.Status == StatusFailed && dm.progress.ErrorMessage == "下载被用户终止" {
				return downloadErr
			}

			// 如果是暂停信号，等待恢复
			if dm.progress.Status == StatusPaused {
				dm.waitForResume()
				existingSize, resumeSupported = dm.checkExistingFile(savePath)
				dm.mutex.Lock()
				dm.progress.Downloaded = existingSize
				dm.mutex.Unlock()
				continue
			}

			// 重试逻辑
			if dm.retries < dm.options.MaxRetries {
				dm.retries++
				DebugWithFormat("下载失败，%d秒后第%d次重试: %v", dm.options.RetryInterval/time.Second, dm.retries, downloadErr)
				time.Sleep(dm.options.RetryInterval)
				existingSize, resumeSupported = dm.checkExistingFile(savePath)
				dm.mutex.Lock()
				dm.progress.Downloaded = existingSize
				dm.mutex.Unlock()
				continue
			}
		}

		return err
	}
}

// 带断点续传的下载实现
func (dm *DownloadManager) downloadWithResume(savePath string, existingSize int64, resumeSupported bool) error {
	// 记录下载开始的调试信息
	DebugWithFormat("开始下载文件，保存路径: %s，%s，大小: %s",
		savePath,
		map[bool]string{true: "支持断点续传", false: "不支持断点续传"}[resumeSupported],
		FormatBytes(existingSize))

	// 创建请求前，先检查URL是否有效（包括重定向处理）
	req, err := http.NewRequest("GET", dm.progress.URL, nil)
	if err != nil {
		return &DownloadError{
			Code:    ErrCodeNetworkError,
			Type:    "请求创建错误",
			Message: "无法创建HTTP请求",
			Err:     err,
		}
	}

	// 设置请求头
	// 使用更现代的User-Agent
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	// 添加接受头部
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	// 对抖音URL添加特殊处理
	if strings.Contains(dm.progress.URL, "aweme.snssdk.com") {
		DebugWithFormat("检测到抖音URL，添加特殊请求头")
		req.Header.Set("Referer", "https://www.iesdouyin.com/")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
	}

	// 如果支持断点续传且已有部分文件，设置Range头
	if resumeSupported && existingSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))
		DebugWithFormat("继续下载，从 %d 字节开始", existingSize)
	}

	// 发送请求
	resp, err := dm.client.Do(req)
	if err != nil {
		// 检查是否是停止信号导致的错误
		select {
		case <-dm.stopChan:
			return &DownloadError{
				Code:    ErrCodeNetworkError,
				Type:    "下载停止",
				Message: "下载被用户终止",
				Err:     err,
			}
		default:
			return &DownloadError{
				Code:    ErrCodeNetworkError,
				Type:    "网络请求错误",
				Message: "发送HTTP请求失败",
				Err:     err,
			}
		}
	}
	defer resp.Body.Close()

	// 记录重定向信息（如果有）
	if resp.Request.URL.String() != dm.progress.URL {
		DebugWithFormat("[redirect] 下载URL重定向: %s -> %s", dm.progress.URL, resp.Request.URL.String())
	}

	// 检查响应状态
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 对于续传请求，206是成功状态码
		if !(resumeSupported && existingSize > 0 && resp.StatusCode == 206) && resp.StatusCode != 200 {
			return &DownloadError{
				Code:    ErrCodeNetworkError,
				Type:    "服务器响应错误",
				Message: fmt.Sprintf("服务器返回错误状态码: %d", resp.StatusCode),
				Err:     fmt.Errorf("HTTP status code: %d", resp.StatusCode),
			}
		}
	}

	// 获取文件总大小
	contentLength := resp.ContentLength
	if contentLength > 0 {
		dm.mutex.Lock()
		// 如果是续传，总大小需要加上已下载部分
		if resumeSupported && existingSize > 0 {
			dm.progress.TotalBytes = existingSize + contentLength
		} else {
			dm.progress.TotalBytes = contentLength
		}
		dm.mutex.Unlock()
	}

	// 打开文件，设置追加模式或创建模式
	fileFlags := os.O_WRONLY | os.O_CREATE
	if resumeSupported && existingSize > 0 {
		fileFlags |= os.O_APPEND
	} else {
		fileFlags |= os.O_TRUNC
	}

	file, err := os.OpenFile(savePath, fileFlags, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return &DownloadError{
				Code:    ErrCodePermissionError,
				Type:    "权限错误",
				Message: "没有文件写入权限",
				Err:     err,
			}
		}
		return &DownloadError{
			Code:    ErrCodeFileWriteError,
			Type:    "文件打开错误",
			Message: "无法打开或创建文件",
			Err:     err,
		}
	}
	defer file.Close()

	// 创建缓冲区读取响应体
	reader := bufio.NewReaderSize(resp.Body, dm.chunkSize)
	buffer := make([]byte, dm.chunkSize)
	lastProgressUpdate := time.Now()

	for {
		// 检查暂停信号
		select {
		case <-dm.pauseChan:
			DebugWithFormat("下载已暂停")
			return nil
		case <-dm.stopChan:
			DebugWithFormat("下载已停止")
			return &DownloadError{
				Code:    ErrCodeNetworkError,
				Type:    "下载停止",
				Message: "下载被用户终止",
				Err:     nil,
			}
		default:
			// 继续下载
		}

		// 读取数据
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return &DownloadError{
				Code:    ErrCodeNetworkError,
				Type:    "数据读取错误",
				Message: "读取响应数据失败",
				Err:     err,
			}
		}

		if n == 0 {
			// 文件读取完毕
			break
		}

		// 写入文件
		written, err := file.Write(buffer[:n])
		if err != nil {
			return &DownloadError{
				Code:    ErrCodeFileWriteError,
				Type:    "文件写入错误",
				Message: "写入文件数据失败",
				Err:     err,
			}
		}

		if written < n {
			return &DownloadError{
				Code:    ErrCodeFileWriteError,
				Type:    "文件写入错误",
				Message: "写入文件数据不完整",
				Err:     fmt.Errorf("只写入了 %d/%d 字节", written, n),
			}
		}

		// 更新进度
		dm.mutex.Lock()
		dm.progress.Downloaded += int64(written)
		currentBytes := dm.progress.Downloaded
		currentTime := time.Now()

		// 计算下载速度
		if currentTime.Sub(dm.startTime) > 1*time.Second && currentBytes > dm.lastBytes {
			elapsed := currentTime.Sub(lastProgressUpdate).Seconds()
			if elapsed > 0 {
				speedBytes := float64(currentBytes-dm.lastBytes) / elapsed
				dm.progress.Speed = speedBytes / 1024                // KB/s
				dm.progress.FormattedSpeed = FormatSpeed(speedBytes) // 使用FormatSpeed
			}
			lastProgressUpdate = currentTime
			dm.lastBytes = currentBytes
		}

		// 计算进度百分比
		if dm.progress.TotalBytes > 0 {
			dm.progress.Progress = float64(currentBytes) / float64(dm.progress.TotalBytes) * 100
		}
		// 更新格式化字段
		dm.progress.FormattedSize = FormatBytes(dm.progress.TotalBytes)
		dm.progress.FormattedDownloaded = FormatBytes(currentBytes)
		dm.progress.LastUpdate = currentTime
		dm.mutex.Unlock()

		// 调用进度回调函数
		if dm.options.ProgressFunc != nil {
			go dm.options.ProgressFunc(dm.GetProgress())
		}
	}

	return nil
}

// 验证下载完整性
func (dm *DownloadManager) verifyDownload(filePath string) error {
	if dm.options.ChecksumType == "" || dm.options.ExpectedHash == "" {
		return nil // 不需要校验
	}

	file, err := os.Open(filePath)
	if err != nil {
		return &DownloadError{
			Code:    ErrCodeFileWriteError,
			Type:    "文件打开错误",
			Message: "无法打开文件进行完整性校验",
			Err:     err,
		}
	}
	defer file.Close()

	var hashValue string
	switch strings.ToLower(dm.options.ChecksumType) {
	case "md5":
		hash := md5.New()
		if _, err := io.Copy(hash, file); err != nil {
			return &DownloadError{
				Code:    ErrCodeChecksumError,
				Type:    "校验和计算错误",
				Message: "计算MD5校验和失败",
				Err:     err,
			}
		}
		hashValue = hex.EncodeToString(hash.Sum(nil))
	case "sha256":
		hash := sha256.New()
		if _, err := io.Copy(hash, file); err != nil {
			return &DownloadError{
				Code:    ErrCodeChecksumError,
				Type:    "校验和计算错误",
				Message: "计算SHA256校验和失败",
				Err:     err,
			}
		}
		hashValue = hex.EncodeToString(hash.Sum(nil))
	default:
		return &DownloadError{
			Code:    ErrCodeInvalidConfig,
			Type:    "配置错误",
			Message: fmt.Sprintf("不支持的校验和类型: %s", dm.options.ChecksumType),
			Err:     nil,
		}
	}

	// 比较校验和
	if strings.ToLower(hashValue) != strings.ToLower(dm.options.ExpectedHash) {
		return &DownloadError{
			Code:    ErrCodeChecksumError,
			Type:    "校验和不匹配",
			Message: fmt.Sprintf("文件完整性校验失败，期望: %s，实际: %s", dm.options.ExpectedHash, hashValue),
			Err:     nil,
		}
	}

	DebugWithFormat("文件完整性校验通过，%s: %s", dm.options.ChecksumType, hashValue)
	return nil
}

// 检查是否存在部分下载的文件
func (dm *DownloadManager) checkExistingFile(savePath string) (int64, bool) {
	fileInfo, err := os.Stat(savePath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, true // 文件不存在，但支持续传
		}
		return 0, false // 其他错误，不支持续传
	}

	if fileInfo.IsDir() {
		return 0, false // 是目录，不支持续传
	}

	return fileInfo.Size(), true
}

// 获取保存路径
func (dm *DownloadManager) getSavePath() (string, error) {
	// 确定基础路径
	basePath := dm.options.SavePath
	if basePath == "" {
		// 使用当前目录作为默认路径
		currentDir, err := os.Getwd()
		if err != nil {
			return "", err
		}
		basePath = currentDir
	}

	// 确保目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", err
	}

	// 确定文件名
	fileName := dm.options.FileName
	if fileName == "" {
		// 从URL中提取文件名
		fileName = extractFileNameFromURL(dm.progress.URL)
		if fileName == "" {
			fileName = fmt.Sprintf("download_%d", time.Now().Unix())
		}
	}

	// 组合完整路径
	savePath := filepath.Join(basePath, fileName)
	DebugWithFormat("文件保存路径: %s", savePath)
	return savePath, nil
}

// 等待恢复下载信号
func (dm *DownloadManager) waitForResume() {
	select {
	case <-dm.resumeChan:
		DebugWithFormat("继续下载")
	case <-dm.stopChan:
		DebugWithFormat("下载被停止")
	}
}

// 创建HTTP客户端
func createHTTPClient(timeout time.Duration, ignoreSSL bool) *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: ignoreSSL,
		},
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 20,
		IdleConnTimeout:     90 * time.Second,
		// 允许重定向（默认行为）
	}

	return &http.Client{
		Transport: transport,
		Timeout:   timeout,
		// 增加重定向次数限制
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("too many redirects")
			}
			return nil
		},
	}
}

// 验证URL合法性
func validateURL(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return err
	}

	// 验证协议
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return errors.New("只支持HTTP和HTTPS协议")
	}

	// 验证域名
	if parsedURL.Host == "" {
		return errors.New("无效的域名")
	}

	// 基本域名格式验证
	hostRegex := regexp.MustCompile(`^([a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+)$`)
	if !hostRegex.MatchString(parsedURL.Host) && parsedURL.Host != "localhost" {
		// 检查是否是IP地址
		ipRegex := regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`)
		if !ipRegex.MatchString(parsedURL.Host) {
			return errors.New("无效的域名或IP地址格式")
		}
	}

	return nil
}

// 从URL提取文件名
func extractFileNameFromURL(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return ""
	}

	// 从路径中提取文件名
	path := parsedURL.Path
	fileName := filepath.Base(path)

	// 如果文件名是空的或者只是扩展名，尝试从查询参数中获取
	if fileName == "" || fileName == "." || fileName == ".." || strings.HasPrefix(fileName, ".") && !strings.Contains(fileName[1:], ".") {
		// 检查常见的文件名参数
		query := parsedURL.Query()
		for _, paramName := range []string{"filename", "file", "name", "fname"} {
			if values := query[paramName]; len(values) > 0 && values[0] != "" {
				return filepath.Base(values[0])
			}
		}
	}

	return fileName
}

// 格式化字节数为可读形式
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB"}
	return fmt.Sprintf("%.1f %s", float64(bytes)/float64(div), units[exp])
}

// 格式化速度
func FormatSpeed(bytesPerSecond float64) string {
	if bytesPerSecond < 1024 {
		return fmt.Sprintf("%.1f B/s", bytesPerSecond)
	} else if bytesPerSecond < 1024*1024 {
		return fmt.Sprintf("%.1f KB/s", bytesPerSecond/1024)
	} else {
		return fmt.Sprintf("%.1f MB/s", bytesPerSecond/(1024*1024))
	}
}
