package config

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 初始化日志。支持 console 和 file 两种输出方式。
func InitLogger() *zap.Logger {
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleWriter := zapcore.Lock(os.Stdout)
	consoleCore := zapcore.NewCore(consoleEncoder, consoleWriter, zap.DebugLevel)

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "app.log",
		MaxSize:    100,  // 最大大小 (MB)
		MaxBackups: 3,    // 最大备份数
		MaxAge:     7,    // 最大保存天数
		Compress:   true, // 是否压缩
		LocalTime:  true, // 使用本地时间
	})
	fileCore := zapcore.NewCore(fileEncoder, fileWriter, zap.InfoLevel)
	core := zapcore.NewTee(consoleCore, fileCore)
	logger := zap.New(core)
	return logger
}
