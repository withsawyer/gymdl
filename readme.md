# 🎵 GYMDL - 跨平台智能音乐下载与管理工具

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)]()
[![License](https://img.shields.io/badge/license-MIT-green)]()
[![Build Status](https://img.shields.io/github/actions/workflow/status/nichuanfang/gymdl/release.yml?logo=github)]()
[![Telegram](https://img.shields.io/badge/telegram-bot-blue?logo=telegram)]()

---

## 项目简介

GYMDL 是一个基于 **Go** 语言开发的现代化音乐下载与管理系统，  
支持多平台音乐链接智能识别、下载、解密、整理及同步，集成了 CookieCloud、WebDAV、Telegram Bot 和 AI 助手等多项功能，助力打造无缝的音乐体验。

---

## 核心功能

- 🎯 **多平台支持**：网易云音乐、Apple Music、Spotify、QQ 音乐、YouTube Music、SoundCloud 等主流音乐平台。
- 🔗 **智能链接解析**：自动识别并处理多种音乐链接格式。
- 🍪 **CookieCloud 集成**：自动同步 Cookie，提升下载成功率。
- ☁️ **WebDAV 支持**：自动整理处理好的音乐文件到本地目录或者通过Webdav上传到NAS 。
- 🤖 **Telegram Bot**：通过 Telegram 机器人实现远程控制与通知。
- ⏰ **定时任务调度**：基于 gocron 实现自动化任务管理。
- 🧠 **AI 助手**：智能问答与辅助功能（规划中）。
- 💻 **WebUI**：未来支持可视化管理界面(规划中)。

---

## 技术栈

- **Go 1.22+**：高性能后端语言。
- **Gin**：轻量级 Web 框架。
- **Zap**：高效结构化日志库。
- **gocron**：灵活的定时任务调度。
- **Telebot**：Telegram Bot 框架。
- **gowebdav**：WebDAV 客户端。
- **OpenAI API**：AI 智能问答支持。
- **GitHub Actions**：CI/CD 自动化构建与发布。

---

## 快速开始

1. 克隆仓库并编译：

   ```bash
   git clone https://github.com/nichuanfang/gymdl.git
   cd gymdl
   make release
   ```

2. 编辑 `config.json`，配置 CookieCloud、WebDAV、Telegram 等参数。

3. 运行程序：

   ```bash
   ./gymdl -c config.json
   ```

4. 通过 Telegram 机器人发送音乐链接，享受自动下载与整理。

---

## 贡献指南

欢迎提交 Issue 和 Pull Request，  
请遵守代码规范，保持代码整洁，详细描述变更内容。

---

## 许可证

本项目采用 [MIT 许可证](LICENSE) 开源。

---

## 联系方式

- GitHub: [nichuanfang](https://github.com/nichuanfang)
- Gmail : f18326186224@gmail.com

---

感谢使用 GYMDL，愿你的音乐永不停歇！ 🎶