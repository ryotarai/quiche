# quiche

```go
import (
    "github.com/ryotarai/quiche"
    qsqlx "github.com/ryotarai/quiche/sqlx"
)

cache = quiche.NewMemory[User]()
// OR
redis, err := rueidis.NewClient(rueidis.ClientOption{
    InitAddress: []string{"127.0.0.1:6379"},
    // CacheSizeEachConn: 128 * (1 << 20), // 128 MiB
})
maxClientTTL := time.Hour // the maximum TTL on the client side
cache = quiche.NewRedis[User](redis, "cache-name", maxClientTTL)

cache.Set("key1", User{})
cache.Get("key1")
cache.Fetch("key1", func() User { return User{} })

// A wrapper of sqlx is also available.
cachedDB := qsqlx.New(db, cache)

cachedDB.Select(db, &people, "SELECT * FROM person ORDER BY first_name ASC")
cachedDB.SelectContext(ctx, db, &people, "SELECT * FROM person ORDER BY first_name ASC")
cachedDB.Get(db, &jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
cachedDB.GetContext(ctx, db, &jason, "SELECT * FROM person WHERE first_name=$1", "Jason")
cachedDB.Invalidate("SELECT * FROM person ORDER BY first_name ASC")
```

## Tuning

### `tracking-table-max-keys` of Redis config

The configuration setting tracking-table-max-keys determines the maximum number of keys stored in the invalidation table and is set to 1000000 keys by default.

### `CacheSizeEachConn` of rueidis.ClientOption

This is the limit size of the local cache.
