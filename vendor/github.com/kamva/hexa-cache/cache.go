package hcache

import (
	"context"
	"errors"
	"time"
)

var ErrKeyNotFound = errors.New("key not found")

type Cache interface {
	Name() string
	Get(ctx context.Context, key string, val interface{}) error
	Set(ctx context.Context, key string, val interface{}) error
	SetWithDefaultTTL(ctx context.Context, key string, val interface{}) error
	SetWithTTL(ctx context.Context, key string, val interface{}, ttl time.Duration) error

	Remove(ctx context.Context, key string) error

	// Purge Clears the cache.
	Purge(ctx context.Context) error
}

type Provider interface {
	Cache(name string) Cache
}
