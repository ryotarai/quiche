package quiche

import (
	"context"
	"testing"
	"time"

	"github.com/redis/rueidis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
		// CacheSizeEachConn: 128 * (1 << 20), // 128 MiB
	})
	require.NoError(t, err)

	m := NewRedis[entry](client, "quiche-test", time.Hour)

	ctx := context.Background()

	m.Delete(ctx, "key")

	assert.NoError(t, m.Delete(ctx, "not-exist-key"))

	_, err = m.Get(ctx, "key")
	assert.Equal(t, ErrNotFound, err)

	v, err := m.Fetch(ctx, "key", func() (entry, error) {
		return entry{
			ID:   123,
			Name: "Alice",
		}, nil
	})
	assert.NoError(t, err)
	assert.Equal(t, 123, v.ID)
	assert.Equal(t, "Alice", v.Name)

	time.Sleep(time.Millisecond * 10) // Wait for the invalidation

	v, err = m.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, 123, v.ID)
	assert.Equal(t, "Alice", v.Name)

	err = m.Set(ctx, "key", entry{
		ID:   456,
		Name: "John",
	})
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 10) // Wait for the invalidation

	v, err = m.Get(ctx, "key")
	assert.NoError(t, err)
	assert.Equal(t, 456, v.ID)
	assert.Equal(t, "John", v.Name)
}
