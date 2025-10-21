package core

import (
	"runtime"
)

// 平台枚举
type Platform int

const (
	Windows Platform = iota // 0
	Linux                   // 1
	Macos                   // 2
	Unknown                 // 3 未知平台
)

// String 方法 — 让 fmt.Println(p) 打印平台名
func (p Platform) String() string {
	switch p {
	case Windows:
		return "Windows"
	case Linux:
		return "Linux"
	case Macos:
		return "MacOS"
	default:
		return "Unknown"
	}
}

// PlatformInfo — 自动检测当前运行的操作系统
func PlatformInfo() Platform {
	switch runtime.GOOS {
	case "windows":
		return Windows
	case "linux":
		return Linux
	case "darwin":
		return Macos
	default:
		return Unknown
	}
}
