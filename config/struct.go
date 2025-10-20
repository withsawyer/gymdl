package config

type Config struct {
	Log    *LogConfig    `json:"log"`    // 日志配置
	Ffmpeg *FfmpegConfig `json:"ffmpeg"` // ffmpeg配置，转码使用
	WebDAV *WebDAVConfig `json:"webdav"` //webdav配置
	AI     *AIConfig     `json:"ai"`     //AI配置
}

type LogConfig struct {
	Mode  int    `json:"mode"`  // 日志模式：1标准输出，2日志文件，3标准输出和日志文件
	Level int    `json:"level"` // 日志等级，0-4分别是：debug，info，warning，error，fatal
	File  string `json:"file"`  // 日志文件路径
}

type FfmpegConfig struct {
	MaxWorker   int    `json:"max_worker"`   // 最大进程数：建议为逻辑CPU个数
	FfmpegPath  string `json:"ffmpeg_path"`  // ffmpeg 可执行文件路径
	FfprobePath string `json:"ffprobe_path"` // ffprobe 可执行文件路径
}

type CookieCloudConfig struct {
	CookieCloudUrl  string `json:"cookiecloud_url"`  //cookiecloud 地址
	CookieCloudUUID string `json:"cookiecloud_uuid"` //cookiecloud uuid
	CookieCloudKEY  string `json:"cookiecloud_key"`  //cookiecloud key
}

type WebDAVConfig struct {
	WebDAVUrl  string `json:"webdav_url"`  //webdav地址
	WebDAVUser string `json:"webdav_user"` //webdav用户名
	WebDAVPass string `json:"webdav_pass"` //webdav密码
	WebDAVDir  string `json:"webdav_dir"`  //wevdav路径
}

type AIConfig struct {
	AIBaseUrl string `json:"ai_base_url"` //baseurl
	AIModel   string `json:"ai_model"`    //使用的模型
	AIApiKey  string `json:"ai_api_key"`  //apiKey
}
