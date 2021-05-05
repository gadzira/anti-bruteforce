package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger struct {
	FileName   string
	LogLevel   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	Compress   bool
}

func New(filename, loglevel string, maxsize, maxbackups, maxage int, localtime bool, compress bool) *Logger {
	return &Logger{
		FileName:   filename,
		LogLevel:   loglevel,
		MaxSize:    maxsize,
		MaxBackups: maxbackups,
		MaxAge:     maxage,
		LocalTime:  localtime,
		Compress:   compress,
	}
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func (l *Logger) InitLogger() *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   l.FileName,
		MaxSize:    l.MaxAge,
		MaxBackups: l.MaxBackups,
		MaxAge:     l.MaxAge,
		LocalTime:  l.LocalTime,
		Compress:   l.Compress,
	}

	// w := zapcore.AddSync(&hook)
	h := zapcore.AddSync(&hook)
	var level zapcore.Level

	switch l.LogLevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "warn":
		level = zap.WarnLevel
	default:
		level = zap.InfoLevel
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = TimeEncoder

	nmws := zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(h))

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		nmws,
		level,
	)

	logger := zap.New(core)
	logger.Info("DefaultLogger init success")

	return logger
}
