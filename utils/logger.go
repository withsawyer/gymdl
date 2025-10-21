package utils

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/nichuanfang/gymdl/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	loggerInstance        *zap.Logger        //性能高
	sugaredLoggerInstance *zap.SugaredLogger //封装了易用的高级方法 但是性能低
)

// 彩色等级输出（加粗对齐）
var colorLevelEncoder = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("\033[1;36mDEBUG\033[0m") // 加粗青色
	case zapcore.InfoLevel:
		enc.AppendString("\033[1;32mINFO \033[0m") // 加粗绿色，补齐宽度
	case zapcore.WarnLevel:
		enc.AppendString("\033[1;33mWARN \033[0m") // 加粗黄色，补齐宽度
	case zapcore.ErrorLevel:
		enc.AppendString("\033[1;31mERROR\033[0m") // 加粗红色
	case zapcore.FatalLevel:
		enc.AppendString("\033[1;35mFATAL\033[0m") // 加粗紫色
	default:
		enc.AppendString(fmt.Sprintf("%-5s", l.CapitalString()))
	}
}

// 对齐字段输出
func paddedCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if caller.Defined {
		enc.AppendString(fmt.Sprintf("%-25s", caller.TrimmedPath())) // 文件:行号左对齐，宽度25
	} else {
		enc.AppendString(fmt.Sprintf("%-25s", ""))
	}
}

// 控制台 Encoder
func newConsoleEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:      "T",
		LevelKey:     "L",
		CallerKey:    "C",
		MessageKey:   "M",
		EncodeTime:   zapcore.TimeEncoderOfLayout("15:04:05"),
		EncodeLevel:  colorLevelEncoder,
		EncodeCaller: paddedCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// 文件 Encoder（无颜色，对齐）
func newFileEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		CallerKey:      "C",
		MessageKey:     "M",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   paddedCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// InitLogger 初始化 Logger
func InitLogger(cfg *config.LogConfig) error {
	var level zapcore.Level
	switch cfg.Level {
	case 1:
		level = zap.DebugLevel
	case 2:
		level = zap.InfoLevel
	case 3:
		level = zap.WarnLevel
	case 4:
		level = zap.ErrorLevel
	default:
		level = zap.FatalLevel
	}
	
	consoleEnc := newConsoleEncoder()
	fileEnc := newFileEncoder()
	
	var cores []zapcore.Core
	
	if cfg.Mode == 1 || cfg.Mode == 3 {
		cores = append(cores, zapcore.NewCore(consoleEnc, zapcore.Lock(os.Stdout), level))
	}
	
	if (cfg.Mode == 2 || cfg.Mode == 3) && cfg.File != "" {
		_ = os.MkdirAll(filepath.Dir(cfg.File), 0755)
		f, err := os.OpenFile(cfg.File, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
		cores = append(cores, zapcore.NewCore(fileEnc, zapcore.AddSync(f), level))
	}
	
	loggerInstance = zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))
	sugaredLoggerInstance = loggerInstance.Sugar()
	return nil
}

//Logger 基础Logger 高性能日志
func Logger() *zap.Logger {
	return loggerInstance
}

// SugaredLogger Logger 增强Logger(快速日志) 用于 开发调试或 快速原型
func SugaredLogger() *zap.SugaredLogger {
	return sugaredLoggerInstance
}

// Success 打印带有✅的成功信息
func Success(args ...interface{}) {
	sugaredLoggerInstance.Infof("✅ %s", fmt.Sprint(args...))
}

// Successf  打印带有✅的格式化成功信息
func Successf(format string, args ...interface{}) {
	sugaredLoggerInstance.Infof("✅ "+format, args...)
}

// Sync 同步日志（用于程序退出前 flush）
func Sync() {
	if loggerInstance != nil {
		_ = loggerInstance.Sync()
	}
}
