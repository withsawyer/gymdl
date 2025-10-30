package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// 🧩 LoadConfig 加载配置并填充默认值
func LoadConfig(file string) *Config {
	// 📁 检查文件是否存在
	_, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Println("⚙️ 配置文件未找到，创建默认配置:", file)
		c := createDefaultConfig()
		saveDefaultConfig(file, c)
		return c
	}
	// 📖 读取文件内容
	bytes, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("❌ 读取配置文件失败:", err)
		c := createDefaultConfig()
		saveDefaultConfig(file, c)
		return c
	}
	// 📄 检查是否为空文件
	if len(bytes) == 0 {
		fmt.Println("⚠️ 配置文件为空，生成默认配置")
		c := createDefaultConfig()
		saveDefaultConfig(file, c)
		return c
	}

	// 🧠 尝试解析YAML
	var c = &Config{}
	err = yaml.Unmarshal(bytes, c)
	if err != nil {
		fmt.Println("⚠️ 配置文件解析失败，创建默认配置:", err)
		backupOldConfig(file, bytes)
		c = createDefaultConfig()
		saveDefaultConfig(file, c)
		return c
	}

	// 🧩 填充缺省值
	c.setDefaults()

	// 📢 打印解析后的配置
	// fmt.Printf("解析后的配置: %+v\n", c)

	return c
}

// 🛠️ createDefaultConfig 创建默认配置
func createDefaultConfig() *Config {
	c := &Config{}
	c.setDefaults()
	return c
}

// 💾 saveDefaultConfig 保存默认配置
func saveDefaultConfig(file string, c *Config) {
	data, err := yaml.Marshal(c)
	if err != nil {
		fmt.Println("❌ 序列化默认配置失败:", err)
		return
	}

	dir := getDir(file)
	if dir != "" {
		_ = os.MkdirAll(dir, 0755)
	}

	err = os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println("❌ 写入默认配置失败:", err)
	} else {
		fmt.Println("✅ 默认配置文件已创建:", file)
	}
}

// 🧾 backupOldConfig 备份旧配置
func backupOldConfig(file string, data []byte) {
	backupFile := file + ".bak"
	err := os.WriteFile(backupFile, data, 0644)
	if err != nil {
		fmt.Println("⚠️ 备份旧配置失败:", err)
	} else {
		fmt.Println("🗂️ 已备份旧配置为:", backupFile)
	}
}

// 📂 getDir 获取文件所在目录
func getDir(path string) string {
	if i := len(path) - 1; i >= 0 {
		for ; i >= 0 && path[i] != '/'; i-- {
		}
		if i > 0 {
			return path[:i]
		}
	}
	return ""
}

func (c *Config) setDefaults() {
	if c.WebConfig == nil {
		c.WebConfig = &WebConfig{Enable: false, AppDomain: "localhost", Https: false, AppPort: 8080, GinMode: "debug"}
	}
	if c.CookieCloud == nil {
		c.CookieCloud = &CookieCloudConfig{
			CookieCloudUrl:  "",
			CookieCloudUUID: "",
			CookieCloudKEY:  "",
			CookieFile:      "cookies.txt",
			CookieFilePath:  "data/temp",
			ExpireTime:      180,
		}
	}
	if c.Tidy == nil {
		c.Tidy = &TidyConfig{Mode: 1, DistDir: "data/dist"}
	}
	if c.WebDAV == nil {
		c.WebDAV = &WebDAVConfig{
			WebDAVUrl:  "",
			WebDAVUser: "",
			WebDAVPass: "",
			WebDAVDir:  "",
		}
	}
	if c.Log == nil {
		c.Log = &LogConfig{Mode: 1, Level: 2, File: "data/logs/run.log"}
	}
	if c.Telegram == nil {
		c.Telegram = &TelegramConfig{Enable: false, Mode: 1}
	}
	if c.AI == nil {
		c.AI = &AIConfig{Enable: false, BaseUrl: "https://api.openai.com/v1", Model: "gpt-3.5-turbo"}
	}
	if c.AdditionalConfig == nil {
		c.AdditionalConfig = &AdditionalConfig{EnableCron: false, EnableDirMonitor: false, MonitorDirs: make([]string, 0)}
	}
	if c.ProxyConfig == nil {
		c.ProxyConfig = &ProxyConfig{
			Enable: false,
		}
	}
}
