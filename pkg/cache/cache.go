package cache

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	ctx    context.Context
	client *redis.Client
}

func New(ctx context.Context) (*Cache, error) {
	cacheAddress := os.Getenv("CACHE_ADDRESS")
	rdb := redis.NewClient(&redis.Options{
		Addr:     cacheAddress,
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Cache{
		ctx:    ctx,
		client: rdb,
	}, nil
}

func (c *Cache) Close() {
	c.client.Close()
}

func (c *Cache) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}

func (c *Cache) Set(key string, value interface{}, expiration time.Duration) {
	c.client.Set(c.ctx, key, value, expiration)
}

func (c *Cache) Del(key string) {
	c.client.Del(c.ctx, key)
}
