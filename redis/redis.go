package quiche

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/rueidis"
	"github.com/ryotarai/quiche"
)

var _ quiche.Cache[int] = &Redis[int]{}

func New[T any](client rueidis.Client, key string, ttl time.Duration) *Redis[T] {
	return &Redis[T]{
		key:    fmt.Sprintf("quiche:%s", key),
		client: client,
		ttl:    ttl,
	}
}

type Redis[T any] struct {
	key    string
	client rueidis.Client
	ttl    time.Duration
}

func (r *Redis[T]) Set(ctx context.Context, key string, value T) error {
	serialized, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Do(ctx, r.client.B().Hset().Key(r.key).FieldValue().FieldValue(key, string(serialized)).Build()).Error()
}

func (r *Redis[T]) Get(ctx context.Context, key string) (T, error) {
	return r.get(ctx, key, true)
}

func (r *Redis[T]) GetWithoutCache(ctx context.Context, key string) (T, error) {
	return r.get(ctx, key, false)
}

func (r *Redis[T]) get(ctx context.Context, key string, withCache bool) (T, error) {
	var result rueidis.RedisResult
	if withCache {
		result = r.client.DoCache(ctx, r.client.B().Hget().Key(r.key).Field(key).Cache(), r.ttl)
	} else {
		result = r.client.Do(ctx, r.client.B().Hget().Key(r.key).Field(key).Build())
	}

	b, err := result.AsBytes()
	if err != nil {
		var zero T
		if rueidis.IsRedisNil(err) {
			return zero, quiche.ErrNotFound
		} else {
			return zero, err
		}
	}

	var ret T
	if err := json.Unmarshal(b, &ret); err != nil {
		var zero T
		return zero, err
	}

	return ret, nil
}

func (r *Redis[T]) GetAll(ctx context.Context) ([]T, error) {
	return r.getall(ctx, true)
}

func (r *Redis[T]) GetAllWithoutCache(ctx context.Context) ([]T, error) {
	return r.getall(ctx, false)
}

func (r *Redis[T]) getall(ctx context.Context, withCache bool) ([]T, error) {
	var result rueidis.RedisResult
	if withCache {
		result = r.client.DoCache(ctx, r.client.B().Hgetall().Key(r.key).Cache(), r.ttl)
	} else {
		result = r.client.Do(ctx, r.client.B().Hgetall().Key(r.key).Build())
	}

	b, err := result.AsBytes()
	if err != nil {
		if rueidis.IsRedisNil(err) {
			return nil, quiche.ErrNotFound
		} else {
			return nil, err
		}
	}

	var ret []T
	if err := json.Unmarshal(b, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Redis[T]) Fetch(ctx context.Context, key string, f func() (T, error)) (T, error) {
	v, err := r.Get(ctx, key)
	if err == nil {
		return v, nil
	} else if err != quiche.ErrNotFound {
		var zero T
		return zero, err
	}

	v, err = f()
	if err != nil {
		var zero T
		return zero, err
	}

	if err := r.Set(ctx, key, v); err != nil {
		var zero T
		return zero, err
	}

	return v, nil
}

func (r *Redis[T]) Delete(ctx context.Context, key string) error {
	return r.client.Do(ctx, r.client.B().Hdel().Key(r.key).Field(key).Build()).Error()
}
