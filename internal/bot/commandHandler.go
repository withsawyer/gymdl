package bot

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"go.uber.org/zap"
	tb "gopkg.in/telebot.v4"
)

// StartCommand 响应 /start 命令
func StartCommand(c tb.Context) error {
	msg := `👋 欢迎来到 GymDL Bot!
我可以帮你管理健身相关的任务和信息 🏋️‍♂️

使用 /help 查看可用命令 📜`

	if err := c.Send(msg, &tb.SendOptions{ParseMode: tb.ModeMarkdown}); err != nil {
		logger.Error("Failed to send start message", zap.Error(err))
		return err
	}
	return nil
}

// HelpCommand 响应 /help 命令，自动生成已注册命令说明
func HelpCommand(c tb.Context) error {
	commands, err := c.Bot().Commands()
	if err != nil {
		logger.Error("Failed to get bot commands", zap.Error(err))
		return c.Send("⚠️ 无法获取命令列表")
	}

	if len(commands) == 0 {
		return c.Send("😅 当前没有可用命令")
	}

	helpMsg := "📖 *可用命令列表:*\n\n"
	for _, cmd := range commands {
		helpMsg += fmt.Sprintf("• `/%s` - %s\n", cmd.Text, cmd.Description)
	}

	helpMsg += "\n✨ *Tip:* 你可以直接在聊天中输入命令来体验功能!"

	if err := c.Send(helpMsg, &tb.SendOptions{ParseMode: tb.ModeMarkdown}); err != nil {
		logger.Error("Failed to send help message", zap.Error(err))
		return err
	}

	return nil
}

func WrapperCommand(c tb.Context) error {
	if !app.cfg.WrapperConfig.Enable {
		return c.Send("请先开启wrapper再使用该指令!")
	}
	if app.cfg.WrapperConfig.AppleId == "" || app.cfg.WrapperConfig.AppleSecret == "" {
		return c.Send("请先配置apple_id和apple_secret再使用该指令!")
	}

	cmd := exec.Command(
		"./wrapper",
		"-L", fmt.Sprintf("%s:%s", app.cfg.WrapperConfig.AppleId, app.cfg.WrapperConfig.AppleSecret),
		"-F",
		"-H", "0.0.0.0",
	)
	cmd.Dir = "/app/wrapper"

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return c.Send(fmt.Sprintf("获取stdout失败: %v", err))
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return c.Send(fmt.Sprintf("获取stderr失败: %v", err))
	}

	if err := cmd.Start(); err != nil {
		return c.Send(fmt.Sprintf("启动wrapper失败: %v", err))
	}
	b := c.Bot()
	sender := c.Sender()
	// 发送初始消息并保存返回的 message，用于后续更新
	msg, _ := b.Send(sender, "wrapper 已启动，日志输出如下：\n")

	var mu sync.Mutex // 保证多goroutine更新同一条消息安全
	var logBuffer strings.Builder

	updateMessage := func(line string) {
		mu.Lock()
		defer mu.Unlock()
		// 限制消息长度，避免Telegram限制
		if logBuffer.Len() > 3500 {
			// 保留最后 3500 个字符
			log := logBuffer.String()
			log = log[len(log)-3500:]
			logBuffer.Reset()
			logBuffer.WriteString(log)
		}
		logBuffer.WriteString(line + "\n")
		// 更新消息
		_, _ = b.Edit(msg, logBuffer.String())
	}

	// 异步读取 stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			updateMessage("stdout: " + scanner.Text())
		}
	}()

	// 异步读取 stderr
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			updateMessage("stderr: " + scanner.Text())
		}
	}()

	// 等待命令退出
	go func() {
		if err := cmd.Wait(); err != nil {
			updateMessage(fmt.Sprintf("wrapper 已退出，错误: %v", err))
		} else {
			updateMessage("wrapper 已正常退出")
		}
	}()

	return nil
}

// WrapperSignInCommand 2FA 签入
func WrapperSignInCommand(c tb.Context) error {
	text := strings.TrimSpace(c.Text())

	// 如果用户直接发送命令但没带参数，提示正确用法
	if text == "" {
		return c.Send("请在命令后输入收到的 6 位 2FA 验证码，例如：\n/signin 123456")
	}

	// 提取第一个6位数字序列
	re := regexp.MustCompile(`\d{6}`)
	code := re.FindString(text)
	if code == "" {
		return c.Send("无效的验证码。请输入6位数字验证码，例如：/signin 123456")
	}

	// 目标文件（相对于你之前的 cmd.Dir = /app/wrapper）
	targetPath := filepath.Join("/app/wrapper", "rootfs", "data", "2fa.txt")

	// 确保目录存在
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return c.Send(fmt.Sprintf("创建目录失败: %v", err))
	}

	// 写文件，使用 0600 权限
	if err := os.WriteFile(targetPath, []byte(code), 0600); err != nil {
		return c.Send(fmt.Sprintf("写入2FA失败: %v", err))
	}

	// 可选：如果需要把文件内容同步到磁盘（避免缓冲），可以调用 Sync（这里用系统默认）
	// f, _ := os.OpenFile(targetPath, os.O_RDWR, 0600)
	// if f != nil { f.Sync(); f.Close() }

	// 回复用户成功
	return c.Send("已接收并写入 2FA 验证码，正在尝试完成登录。")

	//todo 启动wrapper进程 如果认证码不对 提示用户
}

// SetCommands 初始化 Telegram 命令列表
func SetCommands(c tb.Context) error {
	commands := []tb.Command{
		{Text: "start", Description: "启动 Bot 👋"},
		{Text: "help", Description: "获取帮助 📜"},
		{Text: "wrapper", Description: "wrapper认证"},
		{Text: "signin", Description: "wrapper签入"},
	}

	if err := c.Bot().SetCommands(commands); err != nil {
		logger.Error("Failed to set commands", zap.Error(err))
		return err
	}

	const successMsg = "✅ 指令初始化成功，使用 /help 查看可用命令"
	if err := c.Send(successMsg); err != nil {
		logger.Error("Failed to send confirmation message", zap.Error(err))
		return err
	}

	return nil
}
