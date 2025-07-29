package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Client struct{ *redis.Client }

func New(addr string) *Client {
	rdb := redis.NewClient(&redis.Options{Addr: addr})
	if err := rdb.Ping(ctx).Err(); err != nil {
		panic("redis: " + err.Error())
	}
	return &Client{rdb}
}

func (c *Client) SetString(key, val string, ttl time.Duration) {
	_ = c.Set(ctx, key, val, ttl).Err()
}

func (c *Client) GetString(key string) (string, bool) {
	s, err := c.Get(ctx, key).Result()
	return s, err == nil
}
