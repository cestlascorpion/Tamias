package storage

import "context"

type Cache interface {
	GetUri(ctx context.Context, key string) (string, int64, error)
	SetUri(ctx context.Context, key, uri string, ttl int64) error
	Close(ctx context.Context) error
}
