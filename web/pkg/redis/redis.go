package redis

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/skale-5/skalogram/web"
)

type Client struct {
	rc *redis.Client
}

func NewClient(addr string) *Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Client{rc: rc}
}

func (c *Client) CachePost(ctx context.Context, id uuid.UUID, content interface{}, ttl time.Duration) (interface{}, error) {
	err := c.rc.Set(ctx, id.String(), content, ttl).Err()
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (c *Client) GetPost(ctx context.Context, id uuid.UUID) (interface{}, error) {
	val, err := c.rc.Get(ctx, id.String()).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	if err == redis.Nil {
		return nil, web.ErrPostCacheNotFound
	}
	return val, nil
}
