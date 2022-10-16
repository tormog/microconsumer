package cache

import (
	rd "github.com/go-redis/redis"
)

type RedisCache struct {
	*rd.Client
}

func NewCache() RedisCache {
	return RedisCache{
		rd.NewClient(&rd.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		}),
	}
}

func (c *RedisCache) Length(key string) int64 {
	return c.LLen(key).Val()
}

func (c *RedisCache) Push(key string, value interface{}) error {
	if err := c.RPush(key, value).Err(); err != nil {
		return err
	}
	return nil
}

func (c *RedisCache) Pop(key string) (string, error) {
	result, err := c.RPop(key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
