package gormx

import (
	"context"
	"fmt"
	"github.com/firma/framework-common/stores/redisx"
	redisClient "github.com/redis/go-redis/v9"
	"time"
)

type CacheStore struct {
	redis *redisClient.Client
}

var cacheStore *CacheStore

func SetRedis(redisConfig *redisx.Config) Store {
	if cacheStore == nil {
		ctx := context.TODO()
		redis := redisx.MustNew(ctx, redisConfig)
		ping := redis.Get().Ping(ctx)
		if ping.Val() == "" {
			panic(fmt.Errorf("gorm redis cache connect addr %s err: %s", redisConfig.Addr, ping.Err()))
		}
		cacheStore = &CacheStore{redis: redis.Get().Client}
	}
	return cacheStore
}

func (c CacheStore) Set(ctx context.Context, key string, value any, ttl time.Duration) error {

	return c.redis.Set(ctx, key, value, ttl).Err()
}

func (c CacheStore) Get(ctx context.Context, key string) ([]byte, error) {

	return c.redis.Get(ctx, key).Bytes()
}

func (c CacheStore) SaveTagKey(ctx context.Context, tag, key string) error {
	return c.redis.SAdd(ctx, tag, key).Err()
}

func (c CacheStore) RemoveFromTag(ctx context.Context, tag string) error {
	return c.redis.SMembers(ctx, tag).Err()
}
