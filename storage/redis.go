package storage

import (
	"context"
	"time"

	"github.com/cestlascorpion/Tamias/core"
	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(ctx context.Context, config *core.Config) (*Redis, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Network:  config.Redis.Protocol,
		DB:       config.Redis.Database,
		PoolSize: config.Redis.PoolSize,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Errorf("redis ping err %+v", err)
		return nil, err
	}

	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) GetUri(ctx context.Context, key string) (string, int64, error) {
	pipe := r.client.Pipeline()

	getCmd := pipe.Get(ctx, key)
	ttlCmd := pipe.TTL(ctx, key)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return "", 0, err
	}

	uri, err := getCmd.Result()
	if err != nil {
		return "", 0, err
	}

	ttl, err := ttlCmd.Result()
	if err != nil {
		return "", 0, err
	}

	return uri, int64(ttl.Seconds()), nil
}

func (r *Redis) SetUri(ctx context.Context, key, uri string, ttl int64) error {
	_, err := r.client.Set(ctx, key, uri, time.Second*time.Duration(ttl)).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *Redis) Close(ctx context.Context) error {
	return r.client.Close()
}
