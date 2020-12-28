package esutils

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

type DefaultCache struct {
	Cache *redis.Client
}

func NewDefaultCache(cache *redis.Client) *DefaultCache {
	return &DefaultCache{Cache: cache}
}

func (d *DefaultCache) Contains(key string) bool {
	b := d.Cache.Exists(ctx, key).Val()
	if b == 0 {
		return false
	}else {
		return true
	}
}

func (d *DefaultCache) Fetch(key string) (string, error) {
	val,err := d.Cache.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", errors.New(fmt.Sprintf("%s does not exist", key))
	} else if err != nil {
		return "", err
	} else {
		return val,nil
	}
}

func (d *DefaultCache) FetchMulti(keys []string) map[string]string {
	result := make(map[string]string)

	for _, key := range keys {
		if value, err := d.Fetch(key); err == nil {
			result[key] = value
		}
	}

	return result
}

func (d *DefaultCache) Flush() error {
	return d.Cache.FlushAll(ctx).Err()
}

func (d *DefaultCache) Save(key string, value string, lifeTime time.Duration) error {
	return d.Cache.Set(ctx, key, value, lifeTime).Err()
}

func (d *DefaultCache) Delete(key string) error {
	return  d.Cache.Del(ctx, key).Err()
}
