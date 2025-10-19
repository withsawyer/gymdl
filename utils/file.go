package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// IsDir 判断是否为目录
func IsDir(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// DirExist 判断目录是否存在
func DirExist(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsFile 判断是否为文件
func IsFile(file string) bool {
	info, err := os.Stat(file)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// FileExist 文件是否存在
func FileExist(file string) bool {
	if _, err := os.Stat(file); err == nil {
		return true
	}
	return false
}

// CopyFile 拷贝文件
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName) // 这里修正，应该先打开源文件
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()

	// 创建目标目录（如果不存在）
	if err := os.MkdirAll(filepath.Dir(dstName), os.ModePerm); err != nil {
		return 0, err
	}

	dst, err := os.OpenFile(dstName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dst.Close()

	return io.Copy(dst, src)
}

// MoveFile 移动文件（跨分区时会用拷贝+删除）
func MoveFile(dstName, srcName string) (written int64, err error) {
	err = os.Rename(srcName, dstName)
	if err == nil {
		// 同分区直接移动成功
		return 0, nil
	}

	// 如果跨分区，使用拷贝+删除方式
	written, err = CopyFile(dstName, srcName)
	if err != nil {
		return written, err
	}

	err = os.Remove(srcName)
	return written, err
}
