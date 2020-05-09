package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type PremiumKeys struct {
	*pgxpool.Pool
}

func newPremiumKeys(db *pgxpool.Pool) *PremiumKeys {
	return &PremiumKeys{
		db,
	}
}

func (k PremiumKeys) Schema() string {
	return `CREATE TABLE IF NOT EXISTS premium_keys("key" varchar(36) NOT NULL UNIQUE, "length" interval NOT NULL, PRIMARY KEY("key"));`
}

func (k *PremiumKeys) Create(key string, length time.Duration) (err error) {
	_, err = k.Exec(context.Background(), `INSERT INTO premium_keys("key", "length") VALUES($1, $2) ON CONFLICT("key") DO UPDATE SET "length" = $2;`, key, length)
	return
}

func (k *PremiumKeys) Delete(key string) (length time.Duration, e error) {
	if err := k.QueryRow(context.Background(), `DELETE from premium_keys WHERE "key" = $1 RETURNING "expiry";`, key).Scan(&length); err != nil && err != pgx.ErrNoRows {
		e = err
	}

	return
}