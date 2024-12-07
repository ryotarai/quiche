package quiche

import (
	"context"
	"errors"
)

var ErrNotFound = errors.New("not found")

type Cache[T any] interface {
	Set(context.Context, string, T) error

	// Get returns the cached content.
	// This returns ErrNotFound if the key is not found.
	Get(context.Context, string) (T, error)

	Fetch(context.Context, string, func() (T, error)) (T, error)
	Delete(context.Context, string) error
}
