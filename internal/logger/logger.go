package logger

/*
example:

// log one for DEBUG
var w zapcore.WriteSyncer
w = zapcore.AddSync(os.Stdout)

log1 :=  newZapCore(true, zapcore.DebugLevel, w)    // this log level is debug and output is STDOUT

// log two for Error
w := zapcore.AddSync(&lumberjack.Logger{
  Filename:   "/var/log/myapp/foo.log",
  MaxSize:    500, // megabytes
  MaxBackups: 3,
  MaxAge:     28, // days
})
core := zapcore.NewCore(
  zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
  w,
  zap.ErrorLevel,
)
log2 := zap.New(core)

//

	opts := []zap.Option{}
      opts = append(opts, zap.AddCaller())


zlog :=    zap.New(zapcore.NewTee(log1, log2), opts...)

//   just use this, log will send out to file and stdout with define log level
zlog.Info("info")
zlog.Error("error.....")
*/
import (
	"fmt"
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
	fmt.Println("Log file:", l.FileName)
	hook := lumberjack.Logger{
		Filename:   l.FileName,
		MaxSize:    l.MaxAge,
		MaxBackups: l.MaxBackups,
		MaxAge:     l.MaxAge,
		LocalTime:  l.LocalTime,
		Compress:   l.Compress,
	}

	w := zapcore.AddSync(&hook)
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
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		level,
	)
	logger := zap.New(core)
	logger.Info("DefaultLogger init success")

	return logger
}
