# 🎵 GYMDL

### 🚀 跨平台智能音乐下载与管理工具

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)]()
[![License](https://img.shields.io/badge/License-MIT-green)]()
[![Build](https://img.shields.io/github/actions/workflow/status/nichuanfang/gymdl/release.yml?logo=github)]()
[![Telegram](https://img.shields.io/badge/Telegram-Bot-blue?logo=telegram)]()

---

## 🧭 项目简介

**GYMDL** 是基于 Go 的跨平台音乐下载与管理工具，支持多平台智能识别链接、下载、解密、整理，并可同步到 WebDAV、接收 Telegram
通知、使用 AI 助手。

---

## ✨ 核心特性

* 🎯 多平台音乐下载：网易云、Apple Music、Spotify、QQ 音乐、YouTube Music、SoundCloud
* 🔗 智能链接识别与解析
* 🍪 CookieCloud 自动同步登录状态
* ☁️ WebDAV 自动上传整理后的音乐
* 🤖 Telegram Bot 控制下载、接收通知
* ⏰ 定时任务调度（gocron）
* 📂 目录监听（规划中）
* 🧠 AI 助手（规划中）
* 💻 Web UI（规划中）

---

## ⚙️ 快速开始

### 1️⃣ 获取项目并编译

```bash
git clone https://github.com/nichuanfang/gymdl.git
cd gymdl
make release
```

### 2️⃣ 配置文件示例

<details>
<summary>点击展开 YAML 配置示例</summary>

```yaml
# =========================
# Web 服务配置
# =========================
web_config:
  enable: false          # 是否启用 Web 服务
  app_domain: "localhost" # Web 服务域名
  https: false           # 是否启用 HTTPS
  app_port: 9527         # Web 服务端口
  gin_mode: "debug"      # Gin 运行模式: debug/release/test

# =========================
# CookieCloud 配置
# =========================
cookie_cloud:
  cookiecloud_url: ""       # CookieCloud 服务地址
  cookiecloud_uuid: ""      # 用户 UUID
  cookiecloud_key: ""       # 多端需一致
  cookie_file_path: ""      # 本地存储目录
  cookie_file: ""           # Cookie 文件名
  expire_time: 180          # 过期时间（分钟）

# =========================
# 音乐整理配置
# =========================
music_tidy:
  mode: 1               # 整理模式: 1=本地, 2=WebDAV
  dist_dir: "data/dist" # 本地整理目录

# =========================
# WebDAV 配置
# =========================
webdav:
  webdav_url: ""
  webdav_user: ""
  webdav_pass: ""
  webdav_dir: ""

# =========================
# 日志配置
# =========================
log:
  mode: 1
  level: 2
  file: "data/logs/run.log"

# =========================
# Telegram 配置
# =========================
telegram:
  enable: false
  mode: 1
  chat_id: ""
  bot_token: ""
  allowed_users: [ "" ]
  webhook_url: ""
  webhook_port: 9000

# =========================
# AI 配置
# =========================
ai:
  enable: false
  base_url: ""
  model: ""
  api_key: ""
  system_prompt: ""

# =========================
# 附加功能配置
# =========================
additional_config:
  enable_cron: false
  enable_monitor: false
  monitor_dirs: [ "" ]
  enable_wrapper: false

# =========================
# 代理配置
# =========================
proxy:
  enable: false
  scheme: "http"
  host: "127.0.0.1"
  port: 7890
  user: ""
  pass: ""
  auth: false
```

</details>

> ⚡ **提示**：展开查看每个字段的详细注释说明，方便初学者直接修改配置。


### 3️⃣ 运行 GYMDL

```bash
./gymdl -c config.yaml
```

GYMDL 的能力：

* 链接识别与下载
* 音源解密
* 监控下载目录自动解密
* 文件整理并上传到 WebDAV 或本地目录
* Telegram 通知与交互


### 4️⃣ 使用流程

1. 安装 [CookieCloud 插件](https://chrome.google.com/webstore/detail/cookiecloud/ffjiejobkoibkjlhjnlgmcnnigeelbdl)
2. 登录音乐平台并同步 Cookie
3. 配置好 `config.yaml`
4. 通过 Telegram Bot 发送音乐链接，GYMDL 自动处理


### 5️⃣ 高音质下载前置条件

* ✅ 科学上网
* ✅ 登录对应音乐平台账号
* ✅ CookieCloud 已同步
* ✅ 部署方式配置环境：

| 部署方式         | 说明                                                                                                |
|--------------|---------------------------------------------------------------------------------------------------|
| 🐳 Docker 部署 | 仅需配置 `config.yaml` 即可                                                                             |
| 💻 本地部署      | 需额外安装：<br>• Python 3.12+<br>• ffmpeg / ffprobe<br>• N_m3u8DL-RE<br>• MP4Box <br>• wrapper(docker) |

[//]: # (## 📸 示例截图)

[//]: # ()
[//]: # (<details>)

[//]: # (<summary>点击查看示例 Telegram 控制截图</summary>)

[//]: # ()
[//]: # (![Telegram 示例]&#40;https://via.placeholder.com/600x300.png?text=Telegram+Bot+Example&#41;)

[//]: # ()
[//]: # (</details>)

[//]: # ()
[//]: # (<details>)

[//]: # (<summary>点击查看 Web UI / 文件整理截图（规划中）</summary>)

[//]: # ()
[//]: # (![Web UI 示例]&#40;https://via.placeholder.com/600x300.png?text=Web+UI+Example&#41;)

[//]: # ()
[//]: # (</details>)

[//]: # ()
[//]: # (---)

## 🤝 贡献指南

* 提交 **Issue** 或 **Pull Request** ❤️
* 保持代码风格一致
* PR 前使用 `go fmt` 格式化代码
* PR 中详细说明改动内容

## 📜 许可证

MIT License ([LICENSE](LICENSE))

## 📬 联系方式

* GitHub：[@nichuanfang](https://github.com/nichuanfang)
* Email：[f18326186224@gmail.com](mailto:f18326186224@gmail.com)

> 💬 *“愿你的音乐，永不停歇。”* 🎧