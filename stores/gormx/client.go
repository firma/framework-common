package gormx

import (
	"context"
	"github.com/firma/framework-common/stores/redisx"
	"github.com/go-kratos/kratos/v2/log"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	//"time"
)

type (
	Config struct {
		DSN        string
		LogLevel   logger.LogLevel
		RedisCache *redisx.Config
		Logger     log.Logger
	}

	DBManager struct {
		DB *gorm.DB
	}
)

func (m DBManager) GetDB(ctx context.Context) (*gorm.DB, error) {
	return m.DB.WithContext(ctx), nil
}

func MustBuildGormDB(conf Config) *DBManager {
	m, err := BuildDBManager(conf)
	if err != nil {
		panic(err)
	}

	return m
}

func BuildDBManager(conf Config) (*DBManager, error) { //, logger logger.Interface

	mysqlConfig := mysql.Config{
		DSN:                       conf.DSN, // DSN data source name
		DefaultStringSize:         191,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据版本自动配置
	}

	gormLogger := NewGormLogger(conf.Logger)
	gormLogger.LogMode(conf.LogLevel) //// 日志级别
	//gormLogger.conf.SlowThreshold = 300 * time.Millisecond //// 慢 SQL 阈值
	//gormLogger.conf.Colorful = true                        // 禁用彩色打印

	client, err := gorm.Open(
		mysql.New(mysqlConfig), &gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束
			Logger:                                   gormLogger,
		},
	)
	// 缓存 redis 需要实现 注解关闭 后续业务层统一处理 logic->cache->repo
	if conf.RedisCache != nil && len(conf.RedisCache.Addr) > 0 {
		store := SetRedis(conf.RedisCache)
		cacheConfig := &CacheConfig{
			Store:      store, //NewWithDb(redisClient), // OR redis2.New(&redis.Options{Addr:"6379"})
			Serializer: &DefaultJSONSerializer{},
			Prefix:     "db:gorm:",
		}

		cachePlugin := New(cacheConfig)
		err := client.Use(cachePlugin)
		if err != nil {
			return nil, err
		}
	}

	if err != nil {
		return nil, err
	}

	// TODO
	// db.SetMaxOpenConns(conf.MaxOpenConns) // 打开数据库连接的最大数量
	// db.SetMaxIdleConns(conf.MaxIdleConns) // 空闲连接池中连接的最大数量
	// db.SetConnMaxLifetime(conf.ConnMaxLifeTime)

	return &DBManager{
		DB: client,
	}, nil

}
