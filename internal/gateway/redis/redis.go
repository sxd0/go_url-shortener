package redis

import (
	"context"
	"time"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	*goredis.Client
}

type Options struct {
	Addr          string
	DialTimeout   time.Duration
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	PoolSize      int
	MinIdleConns  int
}

func New(ctx context.Context, opts Options) (*Client, error) {
	rdb := goredis.NewClient(&goredis.Options{
		Addr:         opts.Addr,
		DialTimeout:  opts.DialTimeout,
		ReadTimeout:  opts.ReadTimeout,
		WriteTimeout: opts.WriteTimeout,
		PoolSize:     opts.PoolSize,
		MinIdleConns: opts.MinIdleConns,
	})
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return &Client{rdb}, nil
}

func (c *Client) Close() error {
	return c.Client.Close()
}

func (c *Client) SetString(ctx context.Context, key, val string, ttl time.Duration) error {
	return c.Set(ctx, key, val, ttl).Err()
}

func (c *Client) GetString(ctx context.Context, key string) (string, bool, error) {
	s, err := c.Get(ctx, key).Result()
	if err == goredis.Nil {
		return "", false, nil
	}
	if err != nil {
		return "", false, err
	}
	return s, true, nil
}
