package quiche

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemory(t *testing.T) {
	m := NewMemory[entry]()

	ctx := context.Background()

	_, err := m.Get(ctx, "key")
	assert.Equal(t, ErrNotFound, err)

	v, err := m.Fetch(ctx, "key", func() (entry, error) {
		return entry{
			ID:   123,
			Name: "ryotarai",
		}, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 123, v.ID)
	assert.Equal(t, "ryotarai", v.Name)

	v, err = m.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, 123, v.ID)
	assert.Equal(t, "ryotarai", v.Name)
}
