package zap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type LogConfig struct {
	Level             int8   `json:"level"`               // 日志等级
	LogFile           string `json:"log_file"`            // 日志文件
	DisableLogConsole bool   `json:"disable_console_log"` // 是否将日志输出在console
	MaxSize           int    `json:"max_size"`            // MaxSize 进行切割之前，日志文件的最大大小(MB为单位)，默认为100MB
	MaxAge            int    `json:"max_age"`             // MaxAge 是根据文件名中编码的时间戳保留旧日志文件的最大天数。
	MaxBackups        int    `json:"max_backups"`         // MaxBackups 是要保留的旧日志文件的最大数量。默认是保留所有旧的日志文件（尽管 MaxAge 可能仍会导致它们被删除。）
	LocalTime         bool   `json:"local_time"`          // 使用本地时间
	Compress          bool   `json:"compress"`            // 是否压缩
	Skip              int    `json:"skip"`                // 日志打印层级
}

func NewLogger(lc *LogConfig) *zap.Logger {

	//// 调试级别
	debugPriority := zap.LevelEnablerFunc(
		func(lev zapcore.Level) bool {
			return lev == zap.DebugLevel
		},
	)
	// 日志级别
	infoPriority := zap.LevelEnablerFunc(
		func(lev zapcore.Level) bool {
			return lev == zap.InfoLevel
		},
	)
	// 警告级别
	warnPriority := zap.LevelEnablerFunc(
		func(lev zapcore.Level) bool {
			return lev == zap.WarnLevel
		},
	)
	// 错误级别
	errorPriority := zap.LevelEnablerFunc(
		func(lev zapcore.Level) bool {
			return lev >= zap.ErrorLevel
		},
	)

	cores := [...]zapcore.Core{
		zapcore.NewCore(getEncoder(), getLogWriter(lc, fmt.Sprintf("./%s/debug.log", lc.LogFile)), debugPriority),
		zapcore.NewCore(getEncoder(), getLogWriter(lc, fmt.Sprintf("./%s/info.log", lc.LogFile)), infoPriority),
		zapcore.NewCore(getEncoder(), getLogWriter(lc, fmt.Sprintf("./%s/warn.log", lc.LogFile)), warnPriority),
		zapcore.NewCore(getEncoder(), getLogWriter(lc, fmt.Sprintf("./%s/error.log", lc.LogFile)), errorPriority),
	}

	zapLogger := zap.New(zapcore.NewTee(cores[:]...), zap.AddCaller())
	zapLogger = zapLogger.WithOptions(zap.AddCallerSkip(lc.Skip))
	return zapLogger
}

func getLogWriter(lc *LogConfig, fileName string) zapcore.WriteSyncer {
	// 控制台输出
	consoleSyncer := zapcore.AddSync(os.Stdout)
	if lc == nil || lc.LogFile == "" || lc.DisableLogConsole == false {
		// 没配置文件，输出到console
		return zapcore.NewMultiWriteSyncer(consoleSyncer)
	}
	// 文件写入,log文件的配置
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,      // 日志文件位置
		MaxSize:    lc.MaxSize,    // 进行切割之前，日志文件最大值(单位MB)
		MaxBackups: lc.MaxBackups, // 保留旧文件的最大个数
		MaxAge:     lc.MaxAge,     // 最大时间，默认单位 day
		LocalTime:  lc.LocalTime,  // 使用本地时间
		Compress:   lc.Compress,   // 是否压缩
	}
	fileSyncer := zapcore.AddSync(lumberJackLogger)
	// 不在终端打印
	if lc.DisableLogConsole {
		return zapcore.NewMultiWriteSyncer(fileSyncer)
	}
	return zapcore.NewMultiWriteSyncer(consoleSyncer, fileSyncer)
}

func getEncoder() zapcore.Encoder {

	//encoderConfig := zap.NewProductionEncoderConfig()

	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",   // 日志内容key:val， 前面的key设为msg
		LevelKey:   "level", // 日志级别的key设为level
		TimeKey:    "time",
		NameKey:    "logger", // 日志名
		CallerKey:  "caller",
		//EncodeLevel:    zapcore.CapitalLevelEncoder, // CapitalColorLevelEncoder 颜色输出
		EncodeTime:     customTimeEncoder, // 日志时间
		EncodeDuration: zapcore.SecondsDurationEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		//EncodeCaller:     zapcore.FullCallerEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: " ",
	}

	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 时间格式化
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	type appendTimeEncoder interface {
		AppendTimeLayout(time.Time, string)
	}

}
