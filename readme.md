# 🎵 gymdl

> **YouTube Music 高音质下载神器** — 专为 Premium 用户打造，简洁高效，支持多种个性化配置。

![License](https://img.shields.io/badge/license-MIT-green)
![Python](https://img.shields.io/badge/Python-3.8%2B-blue)
![YouTube Music](https://img.shields.io/badge/YouTube%20Music-Premium-red)
![Status](https://img.shields.io/badge/status-active-brightgreen)

---

## ✨ 功能特点

- 🎧 **高音质下载**  
  支持下载 **256kbps AAC** 音质（YouTube Premium 会员专属）。
- 🍪 **简化 Cookie 配置**  
  使用 Chrome 插件 [Get cookies.txt LOCALLY](https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc)，一键导入 Cookie；运行中如检测到过期，会自动跳转至 YouTube Music 提示重新导入。

- ⚙️ **自定义配置文件**  
  灵活控制下载与存储选项，比如下载完成后自动移动到 NAS 音乐待刮削目录。

- 🎯 **最佳音轨智能选择**  
  优先检测并下载 **141 音轨** (AAC)，若不可用则尝试下载 **251 音轨** (Opus) 并使用 `ffmpeg` 转码为 AAC，确保最佳音质。

- 🏷 **元数据补充**  
  自动嵌入标题、上传者、封面、歌词等信息，让音乐文件更完整。

- 🚧 **更多特性开发中...**

---

## ⚠️ 限制与注意事项

> 请适度下载，过度使用可能导致账号封禁，建议使用小号进行测试与下载。

---

## 📦 安装与使用

### 1. 克隆项目

```bash
git clone https://github.com/yourusername/gymdl.git
cd gymdl
```

### 2. 安装依赖

```bash
pip install -r requirements.txt
```

### 3. 配置 Cookie

1. 安装 Chrome 插件 [Get cookies.txt LOCALLY](https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc)
2. 打开 [YouTube Music](https://music.youtube.com/) 并导出 Cookie 到本地文件
3. 在配置文件中填写 Cookie 路径

### 4. 开始下载

```bash
python gymdl.py --config config.yaml
```

---

## 🛠 配置文件示例

```yaml
download_path: './downloads'
move_to_nas: true
nas_path: '/mnt/nas/music_pending'
audio_quality: 'best' # best / aac / opus
embed_metadata: true
cookie_file: './cookies.txt'
```

---

## 📚 技术栈

- **Python 3.8+**
- `yt-dlp` — 视频/音频下载核心
- `ffmpeg` — 音频转码与处理
- 自定义配置解析器 — 灵活的下载控制

---

## 🤝 贡献指南

欢迎提出建议或提交代码改进：

1. Fork 项目
2. 创建新分支：`git checkout -b feature-xxx`
3. 提交更改：`git commit -m 'Add xxx'`
4. 推送分支：`git push origin feature-xxx`
5. 发起 Pull Request

---

## 📄 开源协议

本项目遵循 [MIT License](LICENSE) 开源协议。

---

## 💡 致谢

感谢以下开源工具与库：

- [yt-dlp](https://github.com/yt-dlp/yt-dlp)
- [ffmpeg](https://ffmpeg.org/)
- [Get cookies.txt LOCALLY](https://chromewebstore.google.com/detail/get-cookiestxt-locally/cclelndahbckbenkjhflpdbgdldlbecc)
