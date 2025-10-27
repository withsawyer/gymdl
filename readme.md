<div align="center">

# 🎧 GYMDL — Music Never Sleeps

_一个跨平台音乐下载、CookieCloud深度集成、整理的自动化工具。  
让你的音乐世界，在云端自由流动。_

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Build](https://img.shields.io/github/actions/workflow/status/yourname/cookiesync/build.yml?logo=github)
![Telegram](https://img.shields.io/badge/telegram-bot-blue?logo=telegram)

</div>

---

## 🧠 概念

*GYMDL* 是一个用 **Go** 构建的模块化音乐管理系统，  
集成了 **多平台音乐同步、AI 助手、Telegram 控制、WebDAV 存储** 等功能。

> 🎵 把网易云、Spotify、Apple Music、QQ 音乐、SoundCloud、YouTube Music 全部聚合在一起，  
> 用一个指令，同步、下载、刮削、整理，交给 AI 打理。

---
## 🧩 功能概览

| 模块 | 描述 | 状态 |
|------|------|------|
| ⚙️ 模块化框架 | 支持插件式加载、配置文件驱动 | ✅ 已完成 |
| 🧾 日志系统 | Zap 集成，支持多级别与分模块控制 | ✅ 已完成 |
| ⏰ 定时任务 | 基于 gocron 实现任务调度与中间件 | ✅ 已完成 |
| 🍪 CookieCloud | 自动化同步 Cookie 助力下载 | ✅ 已完成 |
| ☁️ WebDAV 支持 | 上传、同步、备份你的音乐库 | ✅ 已完成 |
| 🤖 Telegram Bot | 控制台与命令入口，带鉴权与日志中间件 | ✅ 已完成 |
| 🎧 音乐平台集成 | 网易云 / Spotify / Apple Music / QQ / YouTube Music / SoundCloud | 🚧 开发中 |
| 🔄 目录监听 | 自动触发刮削、整理、上传流程 | 🚧 规划中 |
| 🧠 AI 模块 | 智能理解你的需求，协助整理与同步 | 🚧 规划中 |
| 💻 WebUI | 可视化界面 | 🚧 规划中 |

---

## 🛠️ 技术栈

```text
Go 1.22+
├── Gin           → Web框架
├── Zap           → 日志系统
├── Gocron        → 定时任务
├── Telebot       → Telegram Bot 框架
├── WebDAV        → 云存储接口
├── Github Action → CI/CD 自动构建
└── AI 模块       → 智能命令与音乐识别
