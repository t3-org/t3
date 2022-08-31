package hcache

import (
	"context"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/tracer"
)

const DELETE_BY_PATTERN_SCRIPT = `
local keys = redis.call('keys', ARGV[1]) 
for i=1,#keys,5000 do
	redis.call('del', unpack(keys, i, math.min(i+4999, #keys)))
end
return keys`

var deleteByPatternScript = redis.NewScript(0, DELETE_BY_PATTERN_SCRIPT)

type redisCache struct {
	// prefix is a global prefix
	prefix     string
	name       string
	pool       *redis.Pool
	marshal    Marshaler
	unmarshal  Unmarshaler
	defaultTTL time.Duration
}

type RedisOptions struct {
	Prefix      string
	Pool        *redis.Pool
	Marshaler   Marshaler
	Unmarshaler Unmarshaler
	DefaultTTL  time.Duration
}

func NewRedisCache(name string, o *RedisOptions) Cache {
	return &redisCache{
		prefix:     o.Prefix,
		name:       name,
		pool:       o.Pool,
		marshal:    o.Marshaler,
		unmarshal:  o.Unmarshaler,
		defaultTTL: o.DefaultTTL,
	}
}

func (c *redisCache) Name() string {
	return c.name
}

func (c *redisCache) key(k string) string {
	// e.g., cache_user_283jf38jf (prefix is "cache_")
	return fmt.Sprintf("%s%s_%s", c.prefix, c.name, k)
}

func (c *redisCache) Get(ctx context.Context, key string, val interface{}) error {
	b, err := redis.Bytes(c.pool.Get().Do("GET", c.key(key)))
	if err != nil {
		if err == redis.ErrNil {
			return ErrKeyNotFound
		}

		return tracer.Trace(err)
	}

	return tracer.Trace(c.unmarshal(b, val))
}

func (c *redisCache) Set(ctx context.Context, key string, val interface{}) error {
	return c.SetWithTTL(ctx, key, val, 0)
}

func (c *redisCache) SetWithDefaultTTL(ctx context.Context, key string, val interface{}) error {
	return c.SetWithTTL(ctx, key, val, c.defaultTTL)
}

func (c *redisCache) SetWithTTL(ctx context.Context, key string, val interface{}, ttl time.Duration) error {
	b, err := c.marshal(val)
	if err != nil {
		return tracer.Trace(err)
	}

	if ttl > 0 {
		_, err = c.pool.Get().Do("SET", c.key(key), b, "px", ttl.Milliseconds())
		return tracer.Trace(err)
	}

	_, err = c.pool.Get().Do("SET", c.key(key), b)
	return tracer.Trace(err)
}

func (c *redisCache) Remove(ctx context.Context, key string) error {
	_, err := c.pool.Get().Do("DEL", c.key(key))
	return tracer.Trace(err)
}

func (c *redisCache) Purge(ctx context.Context) error {
	hlog.Warn("purge cache store", hlog.String("name", c.name), hlog.String("prefix", c.prefix))
	_, err := deleteByPatternScript.Do(c.pool.Get(), c.key("*"))
	return tracer.Trace(err)
}

var _ Cache = &redisCache{}
