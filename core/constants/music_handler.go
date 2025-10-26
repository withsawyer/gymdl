package constants

import "path/filepath"

var BaseTempDir = filepath.Join("data", "temp")

// =============================临时文件夹常量====================================
// 苹果音乐临时文件夹
var AppleMusicTempDir = filepath.Join(BaseTempDir, "AppleMusic")

// 网易云音乐临时文件夹
var NCMTempDir = filepath.Join(BaseTempDir, "NCM")

// QQ音乐临时文件夹
var QQTempDir = filepath.Join(BaseTempDir, "QQ")

// Youtube音乐临时文件夹
var YoutubeTempDir = filepath.Join(BaseTempDir, "Youtube")

// SoundCloud临时文件夹
var SoundcloudTempDir = filepath.Join(BaseTempDir, "Soundcloud")

// Spotify临时文件夹
var SpotifyTempDir = filepath.Join(BaseTempDir, "Spotify")

//====================================================================
