package utils

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// 基本下载示例
func ExampleBasicDownload() {
	fmt.Println("=== 基本下载示例 ===")

	// 下载参数设置
	url := "https://example.com/file.mp3"
	savePath := "./downloads"

	// 创建下载选项
	options := &DownloadOptions{
		SavePath:     savePath,
		FileName:     "file.mp3",
		Timeout:      60 * time.Second, // 60秒超时
		MaxRetries:   3,                // 最多重试3次
		ChecksumType: "",               // 不进行校验和验证
	}

	// 确保下载目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Printf("创建下载目录失败: %v\n", err)
		return
	}

	// 执行下载
	downloader, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动下载
	if err := downloader.Start(); err != nil {
		fmt.Printf("启动下载失败: %v\n", err)
		return
	}

	// 监控下载进度
	go func() {
		lastProgress := 0.0
		for {
			progress := downloader.GetProgress()
			if progress.Status == StatusCompleted || progress.Status == StatusFailed {
				break
			}

			if progress.Progress > lastProgress+5 { // 每5%更新一次
				lastProgress = progress.Progress
				fmt.Printf("下载进度: %.1f%%, 速度: %s/s, 已下载: %s, 总大小: %s\n",
					progress.Progress,
					progress.FormattedSpeed,
					progress.FormattedDownloaded,
					progress.FormattedSize,
				)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	// 等待下载完成
	for {
		progress := downloader.GetProgress()
		if progress.Status == StatusCompleted || progress.Status == StatusFailed {
			if progress.Status == StatusFailed {
				fmt.Printf("下载失败: %s\n", progress.ErrorMessage)
				return
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("下载成功完成")
}

// 大文件下载示例
func ExampleLargeFileDownload() {
	fmt.Println("\n=== 大文件下载示例 ===")

	url := "https://example.com/large-file.zip"
	savePath := "./downloads"

	// 创建下载选项
	options := &DownloadOptions{
		SavePath:   savePath,
		FileName:   "large-file.zip",
		Timeout:    5 * 60 * time.Second, // 5分钟超时
		MaxRetries: 5,                    // 最多重试5次
		ChunkSize:  4 * 1024 * 1024,      // 4MB分块大小
	}

	// 确保下载目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Printf("创建下载目录失败: %v\n", err)
		return
	}

	// 执行下载
	downloader, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动下载
	if err := downloader.Start(); err != nil {
		fmt.Printf("启动下载失败: %v\n", err)
		return
	}

	// 等待下载完成
	for {
		progress := downloader.GetProgress()
		if progress.Status == StatusCompleted || progress.Status == StatusFailed {
			if progress.Status == StatusFailed {
				fmt.Printf("下载失败: %s\n", progress.ErrorMessage)
				return
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("大文件下载成功")
}

// 断点续传示例
func ExampleResumeDownload() {
	fmt.Println("\n=== 断点续传示例 ===")

	url := "https://example.com/very-large-file.iso"

	// 创建下载选项
	options := &DownloadOptions{
		SavePath:     "./downloads",
		FileName:     "very-large-file.iso",
		Timeout:      10 * 60 * time.Second,
		MaxRetries:   10,
		ChecksumType: "sha256", // 使用SHA256验证
	}

	// 执行第一次下载（模拟中途停止）
	fmt.Println("开始第一次下载（模拟部分下载）...")
	downloader1, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动下载
	if err := downloader1.Start(); err != nil {
		fmt.Printf("启动下载失败: %v\n", err)
		return
	}

	// 等待一段时间后停止下载
	time.Sleep(5 * time.Second)
	fmt.Println("停止下载以模拟中断...")
	_ = downloader1.Stop()

	// 等待下载器完全停止
	time.Sleep(1 * time.Second)

	// 执行第二次下载（断点续传）
	fmt.Println("开始断点续传...")
	downloader2, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动续传
	if err := downloader2.Start(); err != nil {
		fmt.Printf("启动续传失败: %v\n", err)
		return
	}

	// 等待下载完成
	for {
		progress := downloader2.GetProgress()
		if progress.Status == StatusCompleted || progress.Status == StatusFailed {
			if progress.Status == StatusFailed {
				fmt.Printf("断点续传失败: %s\n", progress.ErrorMessage)
				return
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("断点续传成功完成")
}

// 并发下载示例
func ExampleConcurrentDownload() {
	fmt.Println("\n=== 并发下载示例 ===")

	urls := []string{
		"https://example.com/file1.mp3",
		"https://example.com/file2.mp3",
		"https://example.com/file3.mp3",
	}

	savePath := "./downloads"
	var wg sync.WaitGroup
	var mutex sync.Mutex
	errors := []error{}

	// 确保下载目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Printf("创建下载目录失败: %v\n", err)
		return
	}

	// 并发下载多个文件
	for i, url := range urls {
		wg.Add(1)
		go func(index int, fileURL string) {
			defer wg.Done()

			options := &DownloadOptions{
				SavePath:   savePath,
				FileName:   fmt.Sprintf("file%d.mp3", index+1),
				Timeout:    60 * time.Second,
				MaxRetries: 3,
			}

			downloader, err := DownloadFile(fileURL, options)
			if err != nil {
				mutex.Lock()
				errors = append(errors, fmt.Errorf("任务 %d 创建失败: %w", index+1, err))
				mutex.Unlock()
				return
			}

			// 启动下载
			if err := downloader.Start(); err != nil {
				mutex.Lock()
				errors = append(errors, fmt.Errorf("任务 %d 启动失败: %w", index+1, err))
				mutex.Unlock()
				return
			}

			// 等待下载完成
			for {
				progress := downloader.GetProgress()
				if progress.Status == StatusCompleted || progress.Status == StatusFailed {
					if progress.Status == StatusFailed {
						mutex.Lock()
						errors = append(errors, fmt.Errorf("任务 %d 失败: %s", index+1, progress.ErrorMessage))
						mutex.Unlock()
					}
					break
				}
				time.Sleep(100 * time.Millisecond)
			}

			fmt.Printf("下载任务 %d 完成\n", index+1)
		}(i, url)
	}

	// 等待所有下载任务完成
	wg.Wait()

	// 打印错误信息
	if len(errors) > 0 {
		fmt.Printf("\n有 %d 个下载任务失败:\n", len(errors))
		for _, err := range errors {
			fmt.Printf("- %v\n", err)
		}
	} else {
		fmt.Println("所有下载任务成功完成")
	}
}

// 自定义请求头示例
func ExampleCustomHeadersDownload() {
	fmt.Println("\n=== 自定义请求头示例 ===")

	url := "https://example.com/protected-file.pdf"
	savePath := "./downloads"

	// 创建下载选项
	options := &DownloadOptions{
		SavePath: savePath,
		FileName: "protected-file.pdf",
		Timeout:  60 * time.Second,
		// 注意：当前实现不直接支持自定义Headers，需要修改底层代码
	}

	// 确保下载目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Printf("创建下载目录失败: %v\n", err)
		return
	}

	// 执行下载
	downloader, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动下载
	if err := downloader.Start(); err != nil {
		fmt.Printf("启动下载失败: %v\n", err)
		return
	}

	// 等待下载完成
	for {
		progress := downloader.GetProgress()
		if progress.Status == StatusCompleted || progress.Status == StatusFailed {
			if progress.Status == StatusFailed {
				fmt.Printf("下载失败: %s\n", progress.ErrorMessage)
				return
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("带自定义配置的下载成功完成")
}

// 可取消下载示例
func ExampleCancellableDownload() {
	fmt.Println("\n=== 可取消下载示例 ===")

	url := "https://example.com/canceled-file.zip"

	// 创建下载选项
	options := &DownloadOptions{
		SavePath: "./downloads",
		FileName: "canceled-file.zip",
		Timeout:  60 * time.Second,
		// 注意：当前实现不直接支持context，需要修改底层代码
	}

	// 执行下载
	downloader, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动下载
	if err := downloader.Start(); err != nil {
		fmt.Printf("启动下载失败: %v\n", err)
		return
	}

	// 设置定时器取消下载
	fmt.Println("开始下载，3秒后取消...")
	time.AfterFunc(3*time.Second, func() {
		fmt.Println("停止下载...")
		_ = downloader.Stop()
	})

	// 等待下载结束
	for i := 0; i < 10; i++ { // 最多等待10秒
		progress := downloader.GetProgress()
		if progress.Status == StatusCompleted || progress.Status == StatusFailed {
			if progress.Status == StatusFailed {
				fmt.Printf("下载结束: %s\n", progress.ErrorMessage)
			} else {
				fmt.Println("下载成功完成（未被取消）")
			}
			break
		}
		time.Sleep(1 * time.Second)
	}
}

// 校验和验证示例
func ExampleChecksumDownload() {
	fmt.Println("\n=== 校验和验证示例 ===")

	url := "https://example.com/file-with-checksum.txt"
	savePath := "./downloads"
	expectedSHA256 := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" // 空文件的SHA256

	// 创建下载选项，启用SHA256校验
	options := &DownloadOptions{
		SavePath:     savePath,
		FileName:     "file-with-checksum.txt",
		ChecksumType: "sha256",
		ExpectedHash: expectedSHA256,
		Timeout:      60 * time.Second,
	}

	// 确保下载目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Printf("创建下载目录失败: %v\n", err)
		return
	}

	// 执行下载
	downloader, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动下载
	if err := downloader.Start(); err != nil {
		fmt.Printf("启动下载失败: %v\n", err)
		return
	}

	// 等待下载完成
	for {
		progress := downloader.GetProgress()
		if progress.Status == StatusCompleted || progress.Status == StatusFailed {
			if progress.Status == StatusFailed {
				fmt.Printf("下载失败: %s\n", progress.ErrorMessage)
				return
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("校验和验证通过，下载成功")
}

// 文件存在检查示例
func ExampleFileExistsDownload() {
	fmt.Println("\n=== 文件存在检查示例 ===")

	url := "https://example.com/example.txt"
	savePath := "./downloads"

	// 创建下载选项
	options := &DownloadOptions{
		SavePath: savePath,
		FileName: "example.txt",
		Timeout:  60 * time.Second,
		// 注意：当前实现没有直接的文件存在处理选项，需要手动检查和处理
	}

	// 手动检查文件是否存在
	fullPath := savePath + "/" + options.FileName
	if _, err := os.Stat(fullPath); err == nil {
		fmt.Printf("文件 %s 已存在\n", fullPath)
		// 根据需要进行处理，如重命名、跳过等
		options.FileName = "example.txt.new"
		fmt.Printf("将使用新文件名: %s\n", options.FileName)
	}

	// 确保下载目录存在
	if err := os.MkdirAll(savePath, 0755); err != nil {
		fmt.Printf("创建下载目录失败: %v\n", err)
		return
	}

	// 执行下载
	downloader, err := DownloadFile(url, options)
	if err != nil {
		fmt.Printf("创建下载器失败: %v\n", err)
		return
	}

	// 启动下载
	if err := downloader.Start(); err != nil {
		fmt.Printf("启动下载失败: %v\n", err)
		return
	}

	// 等待下载完成
	for {
		progress := downloader.GetProgress()
		if progress.Status == StatusCompleted || progress.Status == StatusFailed {
			if progress.Status == StatusFailed {
				fmt.Printf("下载失败: %s\n", progress.ErrorMessage)
				return
			}
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("下载成功完成")
}

// 运行所有下载器示例
func RunDownloaderExamples() {
	fmt.Println("下载器示例运行开始")
	fmt.Println("==================\n")

	// 运行各个示例
	examples := []func(){
		ExampleBasicDownload,
		ExampleLargeFileDownload,
		ExampleResumeDownload,
		ExampleConcurrentDownload,
		ExampleCustomHeadersDownload,
		ExampleCancellableDownload,
		ExampleChecksumDownload,
		ExampleFileExistsDownload,
	}

	for i, example := range examples {
		fmt.Printf("运行示例 %d/%d\n", i+1, len(examples))
		func() {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("示例运行时发生panic: %v\n", r)
				}
			}()
			example()
		}()
		fmt.Println("-------------------")
		time.Sleep(1 * time.Second) // 间隔1秒运行下一个示例
	}

	fmt.Println("\n所有示例运行完成")
}
