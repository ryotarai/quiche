package sqlx

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/ryotarai/quiche"
)

type SqlxDB interface {
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type Wrapper struct {
	db    SqlxDB
	cache quiche.Cache[string]
}

func New(db SqlxDB, cache quiche.Cache[string]) *Wrapper {
	return &Wrapper{
		db:    db,
		cache: cache,
	}
}

func (w *Wrapper) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return w.getOrSelect(w.db.SelectContext, ctx, dest, query, args...)
}

func (w *Wrapper) Select(dest interface{}, query string, args ...interface{}) error {
	return w.SelectContext(context.Background(), dest, query, args...)
}

func (w *Wrapper) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	return w.getOrSelect(w.db.GetContext, ctx, dest, query, args...)
}

func (w *Wrapper) Get(dest interface{}, query string, args ...interface{}) error {
	return w.SelectContext(context.Background(), dest, query, args...)
}

func (w *Wrapper) InvalidateContext(ctx context.Context, query string, args ...interface{}) error {
	key, err := w.cacheKey(query, args)
	if err != nil {
		return err
	}

	return w.cache.Delete(ctx, key)
}

func (w *Wrapper) getOrSelect(f func(ctx context.Context, dest interface{}, query string, args ...interface{}) error, ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	key, err := w.cacheKey(query, args)
	if err != nil {
		return err
	}

	v, err := w.cache.Get(ctx, string(key))
	if err == nil {
		return json.Unmarshal([]byte(v), dest)
	}

	if err := f(ctx, dest, query, args...); err != nil {
		return err
	}

	serialized, err := json.Marshal(dest)
	if err != nil {
		return err
	}

	return w.cache.Set(ctx, string(key), string(serialized))
}

func (w *Wrapper) cacheKey(query string, args interface{}) (string, error) {
	key, err := json.Marshal([]any{query, args})
	if err != nil {
		return "", err
	}
	return string(key), nil
}
