package gormx

import (
	"context"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
)

type CacheConfig struct {
	Store      Store
	Prefix     string
	Serializer Serializer
}

type (
	Serializer interface {
		Serialize(v any) ([]byte, error)

		Deserialize(data []byte, v any) error
	}

	Store interface {
		// Set 写入缓存数据
		Set(ctx context.Context, key string, value any, ttl time.Duration) error

		// Get 获取缓存数据
		Get(ctx context.Context, key string) ([]byte, error)

		// SaveTagKey 将缓存key写入tag
		SaveTagKey(ctx context.Context, tag, key string) error

		// RemoveFromTag 根据缓存tag删除缓存
		RemoveFromTag(ctx context.Context, tag string) error
	}
)

type Cache struct {
	store Store

	// Serializer 序列化
	Serializer Serializer

	// prefix 缓存前缀
	prefix string
}

// New
// @param conf
// @date 2022-07-02 08:09:52
func New(conf *CacheConfig) *Cache {
	if conf.Store == nil {
		os.Exit(1)
	}

	if conf.Serializer == nil {
		conf.Serializer = &DefaultJSONSerializer{}
	}

	return &Cache{
		store:      conf.Store,
		prefix:     conf.Prefix,
		Serializer: conf.Serializer,
	}
}

// Name
// @date 2022-07-02 08:09:48
func (p *Cache) Name() string {
	return "gorm:cache"
}

// Initialize
// @param tx
// @date 2022-07-02 08:09:47
func (p *Cache) Initialize(tx *gorm.DB) error {
	return tx.Callback().Query().Replace("gorm:query", p.Query)
}

// generateKey
// @param key
// @date 2022-07-02 08:09:46
func generateKey(key string) string {
	hash := fnv.New64a()
	_, _ = hash.Write([]byte(key))

	return strconv.FormatUint(hash.Sum64(), 36)
}

// Query
// @param tx
// @date 2022-07-02 08:09:38
func (p *Cache) Query(tx *gorm.DB) {
	ctx := tx.Statement.Context

	var ttl time.Duration
	var hasTTL bool

	if ttl, hasTTL = FromExpiration(ctx); !hasTTL {
		callbacks.Query(tx)
		return
	}

	var (
		key    string
		hasKey bool
	)

	// 调用 Gorm的方法生产SQL
	callbacks.BuildQuerySQL(tx)

	// 是否有自定义key
	if key, hasKey = FromKey(ctx); !hasKey {
		sql, vars := tx.Statement.SQL.String(), tx.Statement.Vars
		sql = strings.Replace(sql, "?", "%v", -1)
		sql = fmt.Sprintf(sql, vars...)
		key = p.prefix + generateKey(sql)
	}

	// 查询缓存数据
	if err := p.QueryCache(ctx, key, tx.Statement.Dest); err == nil {
		return
	}

	// 查询数据库
	p.QueryDB(tx)
	if tx.Error != nil {
		return
	}

	// 写入缓存
	if err := p.SaveCache(ctx, key, tx.Statement.Dest, ttl); err != nil {
		tx.Logger.Error(ctx, err.Error())
		return
	}

	if tag, hasTag := FromTag(ctx); hasTag {
		_ = p.store.SaveTagKey(ctx, tag, key)
	}
}

// QueryDB 查询数据库数据
// 这里重写Query方法 是不想执行 callbacks.BuildQuerySQL 两遍
func (p *Cache) QueryDB(tx *gorm.DB) {
	if tx.Error != nil || tx.DryRun {
		return
	}

	rows, err := tx.Statement.ConnPool.QueryContext(tx.Statement.Context, tx.Statement.SQL.String(), tx.Statement.Vars...)
	if err != nil {
		_ = tx.AddError(err)
		return
	}

	defer func() {
		_ = tx.AddError(rows.Close())
	}()

	gorm.Scan(rows, tx, 0)
}

// QueryCache 查询缓存数据
// @param ctx
// @param key
// @param dest
func (p *Cache) QueryCache(ctx context.Context, key string, dest any) error {

	values, err := p.store.Get(ctx, key)
	if err != nil {
		return err
	}

	switch dest.(type) {
	case *int64:
		dest = 0
	}
	return p.Serializer.Deserialize(values, dest)
}

// SaveCache 写入缓存数据
func (p *Cache) SaveCache(ctx context.Context, key string, dest any, ttl time.Duration) error {
	values, err := p.Serializer.Serialize(dest)
	if err != nil {
		return err
	}

	return p.store.Set(ctx, key, values, ttl)
}

// RemoveFromTag 根据tag删除缓存数据
// @param ctx
// @param tag
// @date 2022-07-02 08:08:59
func (p *Cache) RemoveFromTag(ctx context.Context, tag string) error {
	return p.store.RemoveFromTag(ctx, tag)
}

type (
	// queryCacheCtx
	queryCacheCtx struct{}

	// queryCacheKeyCtx
	queryCacheKeyCtx struct{}

	// queryCacheTagCtx
	queryCacheTagCtx struct{}
)

// NewKey
// @param ctx
// @param key
// @date 2022-07-02 08:11:44
func NewKey(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, queryCacheKeyCtx{}, key)
}

// NewTag
// @param ctx
// @param key
// @date 2022-07-02 08:11:43
func NewTag(ctx context.Context, key string) context.Context {
	return context.WithValue(ctx, queryCacheTagCtx{}, key)
}

// NewExpiration
// @param ctx
// @param ttl
// @date 2022-07-02 08:11:41
func NewExpiration(ctx context.Context, ttl time.Duration) context.Context {
	return context.WithValue(ctx, queryCacheCtx{}, ttl)
}

// FromExpiration
// @param ctx
// @date 2022-07-02 08:11:40
func FromExpiration(ctx context.Context) (time.Duration, bool) {
	value := ctx.Value(queryCacheCtx{})

	if value != nil {
		if t, ok := value.(time.Duration); ok {
			return t, true
		}
	}

	return 0, false
}

// FromKey
// @param ctx
// @date 2022-07-02 08:11:39
func FromKey(ctx context.Context) (string, bool) {
	value := ctx.Value(queryCacheKeyCtx{})

	if value != nil {
		if t, ok := value.(string); ok {
			return t, true
		}

	}

	return "", false
}

// FromTag
// @param ctx
// @date 2022-07-02 08:11:37
func FromTag(ctx context.Context) (string, bool) {
	value := ctx.Value(queryCacheTagCtx{})

	if value != nil {
		if t, ok := value.(string); ok {
			return t, true
		}

	}

	return "", false
}
