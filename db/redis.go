package db

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/redis/go-redis/extra/redisotel/v9"
)

func NewRedis(ctx context.Context, redisOptions redis.Options, options ...Option) (*redis.Client, error) {
	client := redis.NewClient(&redisOptions)
	if _, err := client.Ping(ctx).Result(); err != nil {
		return nil, err
	}

	handler := applyOptions(options...)

	if handler.traces {
		if err := redisotel.InstrumentTracing(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func NewRedisOptions(host, password string) redis.Options {
	return redis.Options{
		Addr:     host,
		Password: password,
	}
}
