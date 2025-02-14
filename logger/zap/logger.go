package zap

import (
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"time"
)

// Config 日志配置
type Config struct {
	Dir        string // 日志根目录
	MaxSize    int    // 每个日志文件最大尺寸，单位 MB
	MaxBackups int    // 保留的旧日志文件最大数量
	MaxAge     int    // 保留的旧日志文件最大天数
	Compress   bool   // 是否压缩旧日志文件
	App        string // 应用名称
	Env        string // 环境名称
}

// ZapLogger 包装 zap.Logger
type ZapLogger struct {
	log  *zap.Logger
	sync bool
}

// getLogWriter 根据级别获取对应的 lumberjack writer
func getLogWriter(config Config, level string) *lumberjack.Logger {
	date := time.Now().Format("2006-01-02")
	filename := filepath.Join(config.Dir, level, fmt.Sprintf("%s.log", date))

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		panic(err)
	}

	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}

// NewZapLogger 创建多级别日志的 ZapLogger
func NewZapLogger(config Config, sync bool) (*ZapLogger, error) {
	// 自定义JSON编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		FunctionKey:   zapcore.OmitKey,
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel: func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(l.String())
		},
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(caller.TrimmedPath())
		},
	}

	// 创建多个 core
	var cores []zapcore.Core

	// 定义级别和对应的 zapcore.Level
	levelMap := map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}

	// 为每个级别创建独立的 core
	for name, level := range levelMap {
		writer := getLogWriter(config, name)
		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(writer),
			level,
		)
		cores = append(cores, core)
	}

	// 控制台输出使用彩色编码器
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(
		consoleEncoder,
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	cores = append(cores, consoleCore)

	// 创建 logger
	core := zapcore.NewTee(cores...)
	logger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(1),
		zap.Fields(
			zap.String("app", config.App),
			zap.String("env", config.Env),
		),
	)

	return &ZapLogger{
		log:  logger,
		sync: sync,
	}, nil
}

// Log 实现 log.Logger 接口
func (l *ZapLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 {
		return nil
	}
	if (len(keyvals) % 2) != 0 {
		keyvals = append(keyvals, "KEYVALS UNPAIRED")
	}

	// 构造 zap 字段
	fields := make([]zap.Field, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		// 处理 key
		key, ok := keyvals[i].(string)
		if !ok {
			key = fmt.Sprint(keyvals[i])
		}
		// 添加其他字段
		fields = append(fields, zap.Any(key, keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug("", fields...)
	case log.LevelInfo:
		l.log.Info("", fields...)
	case log.LevelWarn:
		l.log.Warn("", fields...)
	case log.LevelError:
		l.log.Error("", fields...)
	default:
		l.log.Info("", fields...)
	}

	if l.sync {
		err := l.log.Sync()
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper 创建日志助手
func (l *ZapLogger) Helper() *log.Helper {
	return log.NewHelper(l)
}

// NewLogger 创建日志
func (l *ZapLogger) NewLogger() log.Logger {
	return log.With(l,
		"caller", log.DefaultCaller,
		"traceId", tracing.TraceID(),
		"spanId", tracing.SpanID(),
		"ts", log.DefaultTimestamp,
	)

}
