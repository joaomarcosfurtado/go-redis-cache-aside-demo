package platform

import (
	"context"
	"crypto/tls"
	"time"

	"weather-cache/internal/config"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	if cfg.RedisTLS {
		opt.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewClient(opt)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	return client, nil
}