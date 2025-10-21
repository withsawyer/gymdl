package utils

import (
	"fmt"
	"os"
	"path/filepath"
	
	"github.com/nichuanfang/gymdl/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerInstance *zap.Logger

// 彩色等级输出
var colorLevelEncoder = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	switch l {
	case zapcore.DebugLevel:
		enc.AppendString("\033[36mDEBUG\033[0m") // 青色
	case zapcore.InfoLevel:
		enc.AppendString("\033[32mINFO\033[0m") // 绿色
	case zapcore.WarnLevel:
		enc.AppendString("\033[33mWARN\033[0m") // 黄色
	case zapcore.ErrorLevel:
		enc.AppendString("\033[31mERROR\033[0m") // 红色
	case zapcore.FatalLevel:
		enc.AppendString("\033[35mFATAL\033[0m") // 紫色
	default:
		enc.AppendString(l.CapitalString())
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
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	return zapcore.NewConsoleEncoder(cfg)
}

// 文件 Encoder（无颜色）
func newFileEncoder() zapcore.Encoder {
	cfg := zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		CallerKey:      "C",
		MessageKey:     "M",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return zapcore.NewConsoleEncoder(cfg)
}

//InitLogger 初始化 Logger
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
	return nil
}

//Logger 获取全局 Logger
func Logger() *zap.SugaredLogger {
	if loggerInstance == nil {
		_ = InitLogger(&config.LogConfig{Mode: 1, Level: 2})
	}
	return loggerInstance.Sugar()
}

//Sync 同步日志（用于程序退出前 flush）
func Sync() {
	if loggerInstance != nil {
		_ = loggerInstance.Sync()
	}
}
