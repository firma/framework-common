package gormx

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
	"time"
)

// NewGormLogger 适配 Kratos 日志到 GORM
func NewGormLogger(logger log.Logger) logger.Interface {
	return &GormxCustomerLogger{logger: log.NewHelper(logger)}
}

// GormxCustomerLogger 实现 gorm.Logger 接口
type GormxCustomerLogger struct {
	logger *log.Helper
	level  logger.LogLevel
}

func (l *GormxCustomerLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logger.Log(log.Level(level))
	return l // 可根据需要实现日志级别控制
}

func (l *GormxCustomerLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Info(msg, data)
}

func (l *GormxCustomerLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Warn(msg, data)
}

func (l *GormxCustomerLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Error(msg, data)
}
func (l *GormxCustomerLogger) Errorf(ctx context.Context, msg string, data ...interface{}) {
	l.logger.WithContext(ctx).Errorf(msg, data)
}

func (l *GormxCustomerLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rowsAffected := fc()
	if err != nil {
		l.logger.WithContext(ctx).Errorf(
			"sql error %v, | sql = %v , startTime=%v , rows=%v , elapsed=%v",
			err,
			sql,
			begin.Format(time.RFC3339),
			rowsAffected,
			elapsed,
		)
	} else {
		l.logger.WithContext(ctx).Debugf(
			"client sql = %v time=%v rows=%v elapsed=%v", sql, begin.Format(time.RFC3339), rowsAffected, elapsed,
		)
	}
}
