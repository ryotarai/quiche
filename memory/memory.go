package quiche

import (
	"context"
	"sync"

	"github.com/ryotarai/quiche"
)

var _ quiche.Cache[int] = &Memory[int]{}

func New[T any]() *Memory[T] {
	return &Memory[T]{
		cache: sync.Map{},
	}
}

type Memory[T any] struct {
	cache sync.Map
}

func (m *Memory[T]) Set(ctx context.Context, key string, value T) error {
	m.cache.Store(key, value)
	return nil
}

func (m *Memory[T]) Get(ctx context.Context, key string) (T, error) {
	v, ok := m.cache.Load(key)
	if !ok {
		var zero T
		return zero, quiche.ErrNotFound
	}
	return v.(T), nil
}

func (m *Memory[T]) Fetch(ctx context.Context, key string, f func() (T, error)) (T, error) {
	v, err := m.Get(ctx, key)
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

	if err := m.Set(ctx, key, v); err != nil {
		var zero T
		return zero, err
	}

	return v, nil
}

func (m *Memory[T]) Delete(ctx context.Context, key string) error {
	m.cache.Delete(key)
	return nil
}
