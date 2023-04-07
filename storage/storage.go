package storage

import "context"

type Cache interface {
	GetUri(ctx context.Context, md5, manufacturer string) (string, int64, error)
}
