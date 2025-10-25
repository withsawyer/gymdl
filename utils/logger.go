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
	loggerInstance        *zap.Logger
	sugaredLoggerInstance *zap.SugaredLogger
)

// 日志等级对应的图标
var levelIcons = map[zapcore.Level]string{
	zapcore.DebugLevel:  "🐞",
	zapcore.InfoLevel:   "💡",
	zapcore.WarnLevel:   "⚠️",
	zapcore.ErrorLevel:  "❌",
	zapcore.DPanicLevel: "🔥",
	zapcore.PanicLevel:  "💀",
	zapcore.FatalLevel:  "🛑",
}

// 彩色等级输出（加粗对齐）
var colorLevelEncoder = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("\033[1;36mDEBUG\033[0m")
	case zapcore.InfoLevel:
		enc.AppendString("\033[1;32mINFO \033[0m")
	case zapcore.WarnLevel:
		enc.AppendString("\033[1;33mWARN \033[0m")
	case zapcore.ErrorLevel:
		enc.AppendString("\033[1;31mERROR\033[0m")
	case zapcore.FatalLevel:
		enc.AppendString("\033[1;35mFATAL\033[0m")
	default:
		enc.AppendString(fmt.Sprintf("%-5s", l.CapitalString()))
	}
}

// 对齐字段输出
func paddedCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	if caller.Defined {
		enc.AppendString(fmt.Sprintf("%-25s", caller.TrimmedPath()))
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

// Logger 返回高性能 Logger
func Logger() *zap.Logger {
	return loggerInstance
}

// SugaredLogger 返回快速开发 Logger
func SugaredLogger() *zap.SugaredLogger {
	return sugaredLoggerInstance
}

// Sync 同步日志（程序退出前 flush）
func Sync() {
	if loggerInstance != nil {
		_ = loggerInstance.Sync()
	}
}

// =================== 通用封装 ===================

// logWithIcon 通用函数：带图标输出
func logWithIcon(level zapcore.Level, icon string, msg string) {
	if sugaredLoggerInstance == nil {
		return
	}
	switch level {
	case zapcore.DebugLevel:
		sugaredLoggerInstance.Debugf("%s %s", icon, msg)
	case zapcore.InfoLevel:
		sugaredLoggerInstance.Infof("%s %s", icon, msg)
	case zapcore.WarnLevel:
		sugaredLoggerInstance.Warnf("%s %s", icon, msg)
	case zapcore.ErrorLevel:
		sugaredLoggerInstance.Errorf("%s %s", icon, msg)
	case zapcore.DPanicLevel, zapcore.PanicLevel, zapcore.FatalLevel:
		sugaredLoggerInstance.Errorf("%s %s", icon, msg)
	}
}

// =================== 预定义图标函数 ===================

func Success(args ...interface{}) { logWithIcon(zapcore.InfoLevel, "✅", fmt.Sprint(args...)) }
func Successf(format string, args ...interface{}) {
	logWithIcon(zapcore.InfoLevel, "✅", fmt.Sprintf(format, args...))
}

func ServiceIsOn(args ...interface{}) { logWithIcon(zapcore.InfoLevel, "⚙️", fmt.Sprint(args...)) }
func ServiceIsOnf(format string, args ...interface{}) {
	logWithIcon(zapcore.InfoLevel, "⚙️", fmt.Sprintf(format, args...))
}

func NetworkHealth(args ...interface{}) { logWithIcon(zapcore.InfoLevel, "🌐", fmt.Sprint(args...)) }
func NetworkHealthf(format string, args ...interface{}) {
	logWithIcon(zapcore.InfoLevel, "🌐", fmt.Sprintf(format, args...))
}

func Stop(args ...interface{}) { logWithIcon(zapcore.InfoLevel, "🛑", fmt.Sprint(args...)) }
func Stopf(format string, args ...interface{}) {
	logWithIcon(zapcore.InfoLevel, "🛑", fmt.Sprintf(format, args...))
}

func Info(args ...interface{}) { logWithIcon(zapcore.InfoLevel, "💡", fmt.Sprint(args...)) }
func Infof(format string, args ...interface{}) {
	logWithIcon(zapcore.InfoLevel, "💡", fmt.Sprintf(format, args...))
}

func Warning(args ...interface{}) { logWithIcon(zapcore.WarnLevel, "⚠️", fmt.Sprint(args...)) }
func Warningf(format string, args ...interface{}) {
	logWithIcon(zapcore.WarnLevel, "⚠️", fmt.Sprintf(format, args...))
}

func Debug(args ...interface{}) { logWithIcon(zapcore.DebugLevel, "🐞", fmt.Sprint(args...)) }
func Debugf(format string, args ...interface{}) {
	logWithIcon(zapcore.DebugLevel, "🐞", fmt.Sprintf(format, args...))
}

func Critical(args ...interface{}) { logWithIcon(zapcore.ErrorLevel, "🔥", fmt.Sprint(args...)) }
func Criticalf(format string, args ...interface{}) {
	logWithIcon(zapcore.ErrorLevel, "🔥", fmt.Sprintf(format, args...))
}
