package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/nichuanfang/gymdl/config"
	"github.com/studio-b12/gowebdav"
)

type WebDAV struct {
	Config          *config.WebDAVConfig
	Client          *gowebdav.Client
	lastCheck       time.Time
	lastCheckResult bool
	checkMutex      sync.Mutex
}

var (
	GlobalWebDAV *WebDAV
)

// InitWebDAV 初始化全局 WebDAV，只会执行一次
func InitWebDAV(cfg *config.WebDAVConfig) {
	if cfg == nil || cfg.WebDAVUrl == "" || cfg.WebDAVUser == "" || cfg.WebDAVPass == "" {
		panic("⚠️ WebDAV config is invalid")
	}

	client := gowebdav.NewClient(cfg.WebDAVUrl, cfg.WebDAVUser, cfg.WebDAVPass)
	if err := client.Connect(); err != nil {
		panic(fmt.Sprintf("⚠️ Failed to connect WebDAV: %v", err))
	}

	GlobalWebDAV = &WebDAV{
		Config: cfg,
		Client: client,
	}
}

// -------------------- 连接检测 --------------------

func (w *WebDAV) CheckConnection() bool {
	w.checkMutex.Lock()
	defer w.checkMutex.Unlock()

	if time.Since(w.lastCheck) < time.Minute {
		return w.lastCheckResult
	}

	err := w.Client.Connect()
	w.lastCheck = time.Now()
	w.lastCheckResult = err == nil

	if err != nil {
		logger.Warn(fmt.Sprintf("⚠️ WebDAV connection check failed: %v", err))
	}

	return w.lastCheckResult
}

// -------------------- 文件操作 --------------------

// Upload 上传文件到根目录
func (w *WebDAV) Upload(localPath string) error {
	if localPath == "" {
		return fmt.Errorf("localPath cannot be empty")
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer file.Close()

	// 获取本地文件名
	fileName := filepath.Base(localPath)
	// 构造远程路径（根目录下）
	remoteFullPath := "/" + fileName
	logger.Info("💡start uploading file to webdav...")
	if err := w.Client.WriteStream(remoteFullPath, file, 0644); err != nil {
		logger.Warn(fmt.Sprintf("WebDAV upload failed for %s: %v", remoteFullPath, err))
		return err
	}
	logger.Info(fmt.Sprintf("💡 WebDAV uploaded file successfully: %s", remoteFullPath))
	return nil
}

// UploadTo 上传到指定目录
func (w *WebDAV) UploadTo(localPath, remoteDir string) error {
	if localPath == "" {
		return fmt.Errorf("localPath cannot be empty")
	}
	if remoteDir == "" {
		remoteDir = "/" // 默认根目录
	}

	file, err := os.Open(localPath)
	if err != nil {
		return fmt.Errorf("failed to open local file: %v", err)
	}
	defer file.Close()

	fileName := filepath.Base(localPath)

	// 规范化 remoteDir，确保以 / 开头，不以 / 结尾
	remoteDir = filepath.ToSlash(remoteDir)
	if !strings.HasPrefix(remoteDir, "/") {
		remoteDir = "/" + remoteDir
	}
	remoteDir = strings.TrimRight(remoteDir, "/")

	remoteFullPath := remoteDir + "/" + fileName

	// 确保远程目录存在（如果有 ensureRemoteDir 方法）
	if err := w.ensureRemoteDir(remoteDir); err != nil {
		return fmt.Errorf("failed to ensure remote dir: %v", err)
	}
	logger.Info("💡start uploading file to webdav...")
	if err := w.Client.WriteStream(remoteFullPath, file, 0644); err != nil {
		logger.Warn(fmt.Sprintf("WebDAV upload failed for %s: %v", remoteFullPath, err))
		return err
	}

	logger.Info(fmt.Sprintf("💡 WebDAV uploaded file successfully: %s", remoteFullPath))
	return nil
}

func (w *WebDAV) Download(remotePath, localPath string) error {
	if remotePath == "" || localPath == "" {
		return fmt.Errorf("remotePath and localPath cannot be empty")
	}

	remoteFullPath := w.makeRemotePath(remotePath)
	data, err := w.Client.Read(remoteFullPath)
	if err != nil {
		logger.Warn(fmt.Sprintf("⚠️ WebDAV download failed for %s: %v", remotePath, err))
		return fmt.Errorf("failed to read remote file: %v", err)
	}

	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		return fmt.Errorf("failed to create local directories: %v", err)
	}
	logger.Info("start uploading file to webdav...")
	if err := os.WriteFile(localPath, data, 0644); err != nil {
		logger.Warn(fmt.Sprintf("⚠️ WebDAV failed to write local file: %v", err))
		return err
	}

	logger.Info(fmt.Sprintf("💡 WebDAV downloaded file successfully: %s", remotePath))
	return nil
}

func (w *WebDAV) Delete(remotePath string) error {
	if remotePath == "" {
		return fmt.Errorf("remotePath cannot be empty")
	}

	remoteFullPath := w.makeRemotePath(remotePath)
	if err := w.Client.Remove(remoteFullPath); err != nil {
		logger.Warn(fmt.Sprintf("⚠️ WebDAV delete failed for %s: %v", remotePath, err))
		return err
	}

	logger.Info(fmt.Sprintf("💡 WebDAV deleted file successfully: %s", remotePath))
	return nil
}

func (w *WebDAV) List(remoteDir string) ([]string, error) {
	dir := w.makeRemotePath(remoteDir)
	files, err := w.Client.ReadDir(dir)
	if err != nil {
		logger.Warn(fmt.Sprintf("⚠️ WebDAV list failed for %s: %v", remoteDir, err))
		return nil, err
	}

	var names []string
	for _, f := range files {
		names = append(names, f.Name())
	}
	return names, nil
}

// -------------------- 工具方法 --------------------

func (w *WebDAV) makeRemotePath(path string) string {
	return strings.TrimRight(w.Config.WebDAVDir, "/") + "/" + strings.TrimLeft(path, "/")
}

func (w *WebDAV) ensureRemoteDir(dir string) error {
	if err := w.Client.MkdirAll(dir, 0755); err != nil {
		logger.Warn(fmt.Sprintf("⚠️ WebDAV failed to create remote directory %s: %v", dir, err))
		return err
	}
	return nil
}

// -------------------- 可选参数 --------------------

func WithDir(dir string) func(*config.WebDAVConfig) {
	return func(cfg *config.WebDAVConfig) {
		cfg.WebDAVDir = dir
	}
}
