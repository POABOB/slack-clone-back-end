package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Log 全局日誌實例
	Log *zap.Logger
)

// InitLogger 初始化日誌系統
func InitLogger(level string) error {
	// 設置日誌級別
	var logLevel zapcore.Level
	switch level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	// 配置編碼器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 創建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(zapcore.NewConsoleEncoder(encoderConfig)),
		logLevel,
	)

	// 創建日誌實例
	Log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return nil
}

// Debug 輸出調試日誌
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Info 輸出信息日誌
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn 輸出警告日誌
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error 輸出錯誤日誌
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Fatal 輸出致命錯誤日誌
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}
