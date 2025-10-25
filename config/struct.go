package config

type Config struct {
	WebConfig        *WebConfig         `json:"web_config"`        // web配置
	CookieCloud      *CookieCloudConfig `json:"cookie_cloud"`      // cookiecloud配置
	MusicTidy        *MusicTidyConfig   `json:"music_tidy"`        // 音乐整理配置
	WebDAV           *WebDAVConfig      `json:"webdav"`            // webdav配置
	Log              *LogConfig         `json:"log"`               // 日志配置
	Telegram         TelegramConfig     `json:"telegram"`          // telegram配置
	AI               *AIConfig          `json:"ai"`                // AI配置
	AdditionalConfig *AdditionalConfig  `json:"additional_config"` //附属配置
}

type WebConfig struct {
	Enable    bool   `json:"enable"`     //是否启用该服务
	AppDomain string `json:"app_domain"` // web服务domain
	Https     bool   `json:"https"`      // 是否开启了 https
	AppPort   int    `json:"app_port"`   // web服务监听端口
	GinMode   string `json:"gin_mode"`   // Gin的运行模式: 可选项[debug release test]
}

type CookieCloudConfig struct {
	CookieCloudUrl  string   `json:"cookiecloud_url"`  // cookiecloud 地址
	CookieCloudUUID string   `json:"cookiecloud_uuid"` // cookiecloud uuid
	CookieCloudKEY  []string `json:"cookiecloud_key"`  // cookiecloud key 数组
	CookieFile      string   `json:"cookie_file"`      // cookie文件名
	CookieFilePath  string   `json:"cookie_file_path"` // cookie存储目录
	ExpireTime      int      `json:"expire_time"`      // cookie文件过期时间(分钟) 根据自己需求控制
}

type MusicTidyConfig struct {
	Mode    int    `json:"mode"`     // 音乐整理模式: 1整理到DistDir, 2整理到webdav目录WebDAVDir
	DistDir string `json:"dist_dir"` // 整理到的路径,仅在Mode为1时使用该目录
}

type WebDAVConfig struct {
	WebDAVUrl  string `json:"webdav_url"`  // webdav地址
	WebDAVUser string `json:"webdav_user"` // webdav用户名
	WebDAVPass string `json:"webdav_pass"` // webdav密码
	WebDAVDir  string `json:"webdav_dir"`  // wevdav路径
}

type LogConfig struct {
	Mode  int    `json:"mode"`  // 日志模式：1标准输出，2日志文件，3标准输出和日志文件
	Level int    `json:"level"` // 日志等级，1=debug 2=info 3=warn 4=error 5=fatal
	File  string `json:"file"`  // 日志文件路径
}

type TelegramConfig struct {
	Enable       bool     `json:"enable"`        //是否启用该服务
	Mode         int      `json:"mode"`          // 运行模式: 1长轮询, 2Webhook   开发环境推荐1,生产环境推荐2
	ChatID       string   `json:"chat_id"`       // chat_id 机器人ID
	BotToken     string   `json:"bot_token"`     // telegram机器人token
	AllowedUsers []string `json:"allowed_users"` //白名单列表 填用户ID
	WebhookURL   string   `json:"webhook_url"`   //webhook地址 运行模式为2时必填
	WebhookPort  int      `json:"webhook_port"`  //webhook模式监听的端口
}

type AIConfig struct {
	Enable       bool   `json:"enable"`        //是否开启AI
	BaseUrl      string `json:"base_url"`      // baseurl
	Model        string `json:"model"`         // 使用的模型
	ApiKey       string `json:"api_key"`       // apiKey
	SystemPrompt string `json:"system_prompt"` // 默认系统提示词
}

type AdditionalConfig struct {
	EnableCron bool `json:"enable_cron"` //是否启用定时任务
}
