package quiche

import (
	"context"
	"testing"

	"github.com/ryotarai/quiche"
	"github.com/stretchr/testify/assert"
)

type entry struct {
	ID   int
	Name string
}

func TestMemory(t *testing.T) {
	m := New[entry]()

	ctx := context.Background()

	_, err := m.Get(ctx, "key")
	assert.Equal(t, quiche.ErrNotFound, err)

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
