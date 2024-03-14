package redis

import (
	"context"
	"fmt"
	"github.com/Verce11o/Hezzl-Go/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(ctx context.Context, cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Name,
		Username: cfg.Redis.User,
	})

	err := client.Ping(ctx).Err()

	if err != nil {
		panic(err)
	}

	return client
}
