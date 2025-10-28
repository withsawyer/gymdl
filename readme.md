# 🎵 **GYMDL**

### 🚀 跨平台智能音乐下载与管理工具

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)]()
[![License](https://img.shields.io/badge/License-MIT-green)]()
[![Build](https://img.shields.io/github/actions/workflow/status/nichuanfang/gymdl/release.yml?logo=github)]()
[![Telegram](https://img.shields.io/badge/Telegram-Bot-blue?logo=telegram)]()

---

## 🧭 项目简介

**GYMDL** 是一款基于 **Go** 语言构建的现代化音乐下载与管理系统。
它能智能识别来自多平台的音乐链接，完成 **下载、解密、整理与同步**，并集成了 **CookieCloud、WebDAV、Telegram Bot、AI 助手** 等多种功能，
为你打造一套无缝衔接的全平台音乐体验。 🎶

---

## ✨ 核心特性

* 🎯 **多平台支持**：支持网易云、Apple Music、Spotify、QQ 音乐、YouTube Music、SoundCloud 等主流平台。
* 🔗 **智能链接识别**：自动识别不同格式的音乐链接并解析下载。
* 🍪 **CookieCloud 集成**：自动同步登录状态，提升下载成功率与音质。
* ☁️ **WebDAV 支持**：可将处理后的音乐自动上传至 NAS。
* 🤖 **Telegram Bot 控制**：通过机器人远程控制下载与接收通知。
* ⏰ **定时任务调度**：内置基于 `gocron` 的自动化任务系统。
* 📂 **目录监听（规划中）**：监听网易云/QQ下载目录,实时解锁整理。
* 🧠 **AI 助手（规划中）**：提供音乐相关问答与智能辅助。
* 💻 **Web UI（规划中）**：即将支持可视化的管理界面。

---

## 🧩 技术栈

| 模块           | 技术            |
|--------------| ------------- |
| 核心语言         | Go 1.24       |
| Web 框架       | Gin           |
| 日志系统         | Zap           |
| 定时调度         | gocron        |
| 目录监控         | fsnotify        |
| Telegram Bot | Telebot       |
| 云存储          | gowebdav      |
| AI 支持        | OpenAI API    |
| CI/CD        | GitHub Actions |

---

## ⚙️ 快速开始

### 1️⃣ 获取项目并编译

```bash
git clone https://github.com/nichuanfang/gymdl.git
cd gymdl
make release
```

### 2️⃣ 配置文件

编辑 `config.json`，填入：

* CookieCloud 云端信息
* WebDAV 存储配置
* Telegram Bot Token

### 3️⃣ 运行

```bash
./gymdl -c config.json
```

### 4️⃣ 使用

通过 **Telegram 机器人** 发送音乐链接，GYMDL 将自动：

* 识别链接来源
* 下载并解密音源
* 整理命名并上传至 WebDAV 或本地目录

---

## 📘 开发教程

1. 安装 **Go 1.24+**
2. 克隆仓库
3. 安装依赖

   ```bash
   go mod tidy
   ```
4. 启动本地开发环境

   ```bash
   go run main.go -c config.json
   ```
5. 修改代码后，可通过 `make release` 构建二进制文件。

---

## 💡 使用教程

1. 在浏览器中安装 [CookieCloud 插件](https://chromewebstore.google.com/detail/cookiecloud/ffjiejobkoibkjlhjnlgmcnnigeelbdl?hl=en)
2. 登录各音乐平台并同步 Cookie 到云端
3. 配置好 `config.json`
4. 运行 GYMDL 后，通过 Telegram 机器人控制下载或查看状态

---

## 🔧 前置条件

> 若想下载高音质音乐，请确保以下条件满足：

* ✅ 科学上网问题请自行解决
* ✅ 已登录对应音乐平台账号
* ✅ 使用 [CookieCloud 插件](https://chromewebstore.google.com/detail/cookiecloud/ffjiejobkoibkjlhjnlgmcnnigeelbdl?hl=en) 同步 Cookie
* ✅ 根据部署方式配置环境：

| 部署方式             | 要求                                                                                              |
|------------------|-------------------------------------------------------------------------------------------------|
| 🐳 **Docker 部署** | 仅需正确配置 `config.json`                                                                            |
| 💻 **本地部署**      | 除配置文件外，还需配置以下环境才可以解锁全部服务：<br>• 已安装 **python3.12+**<br>• 已安装 **N_m3u8DL-RE**<br>• 已安装 **um cli** |
---

## 🤝 贡献指南

欢迎提交 **Issue** 与 **Pull Request** ❤️
请遵循以下原则：

* 保持代码风格一致
* 提交前通过 `go fmt` 格式化代码
* 在 PR 中详细说明变更内容

---

## 📜 许可证

本项目使用 [MIT License](LICENSE) 开源。

---

## 📬 联系方式

* **GitHub**：[@nichuanfang](https://github.com/nichuanfang)
* **Email**：[f18326186224@gmail.com](mailto:f18326186224@gmail.com)

---

> 💬 *“愿你的音乐，永不停歇。”* 🎧