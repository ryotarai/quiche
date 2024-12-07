package sqlx

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/redis/rueidis"
	"github.com/ryotarai/quiche"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Person struct {
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
	Email     string
}

type Place struct {
	Country string
	City    sql.NullString
	TelCode int
}

func TestSelect(t *testing.T) {
	driver, mock, err := sqlmock.New()
	require.NoError(t, err)

	db := sqlx.NewDb(driver, "sqlmock")

	client, err := rueidis.NewClient(rueidis.ClientOption{
		InitAddress: []string{"127.0.0.1:6379"},
		// CacheSizeEachConn: 128 * (1 << 20), // 128 MiB
	})
	require.NoError(t, err)

	cache := quiche.NewRedis[string](client, "quiche-sqlx-test", time.Hour)
	cachedDB := New(db, cache)

	mock.MatchExpectationsInOrder(true)
	mock.ExpectQuery(`SELECT \* FROM place ORDER BY telcode ASC`).WillReturnRows(sqlmock.NewRows([]string{"country", "city", "telcode"}).AddRow("US", "New York", 1).AddRow("JP", "Tokyo", 81))

	var places []Place
	err = cachedDB.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
	assert.NoError(t, err)
	t.Log(places)

	places = []Place{}
	err = cachedDB.Select(&places, "SELECT * FROM place ORDER BY telcode ASC")
	assert.NoError(t, err)
	t.Log(places)

	assert.NoError(t, mock.ExpectationsWereMet())
}
