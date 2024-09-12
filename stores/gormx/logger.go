package gormx

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
	"time"
)

type Interface interface {
	LogMode(level logger.LogLevel) Interface
	Info(context.Context, string, ...interface{})
	Warn(context.Context, string, ...interface{})
	Error(context.Context, string, ...interface{})
	Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
}

type GormLogger struct {
	Helper *log.Helper
	conf   logger.Config
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.conf.LogLevel = level
	return l
}

func NewGormLogger(l *log.Helper) *GormLogger {
	return &GormLogger{
		Helper: l,
		conf:   logger.Config{SlowThreshold: 200 * time.Millisecond},
	}

}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Helper.WithContext(ctx).Infof(msg, data)
}
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Helper.WithContext(ctx).Errorf(msg, data)
}
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Helper.WithContext(ctx).Errorf(msg, data)
}

func (l *GormLogger) withContext(ctx context.Context) context.Context {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.Request.Context()
	} else {
		return ctx
	}
}

func (l *GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	elapsed := time.Since(begin)
	ctx = l.withContext(ctx)
	sql, rowsAffected := fc()
	if err != nil {
		l.Helper.WithContext(ctx).Errorf(
			"Sql Error %v, | sql = %v time=%v rows=%v elapsed=%v",
			err,
			sql,
			begin.Format(time.RFC3339),
			rowsAffected,
			elapsed,
		)
	} else if l.conf.LogLevel == logger.Info {
		l.Helper.WithContext(ctx).Infof(
			"Sql Info, | sql = %v time=%v rows=%v elapsed=%v", sql, begin.Format(time.RFC3339), rowsAffected, elapsed,
		)
	} else if l.conf.LogLevel == logger.Silent {
		l.Helper.WithContext(ctx).Infof(
			"Sql Silent, | sql = %v time=%v rows=%v elapsed=%v", sql, begin.Format(time.RFC3339), rowsAffected, elapsed,
		)
	} else if l.conf.LogLevel == logger.Warn {
		l.Helper.WithContext(ctx).Warnf(
			"Sql Warn, | sql = %v time=%v rows=%v elapsed=%v", sql, begin.Format(time.RFC3339), rowsAffected, elapsed,
		)
	} else if l.conf.SlowThreshold != 0 && elapsed > l.conf.SlowThreshold {
		l.Helper.WithContext(ctx).Warnf(
			"Sql Slow, | sql = %v time=%v rows=%v elapsed=%v", sql, begin.Format(time.RFC3339), rowsAffected, elapsed,
		)
	}

}
