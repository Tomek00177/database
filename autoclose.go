package database

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type AutoCloseTable struct {
	*pgxpool.Pool
}

type AutoCloseSettings struct {
	Enabled                 bool
	SinceOpenWithNoResponse *time.Time
	SinceLastMessage        *time.Time
}

func newAutoCloseTable(db *pgxpool.Pool) *AutoCloseTable {
	return &AutoCloseTable{
		db,
	}
}

func (a AutoCloseTable) Schema() string {
	return `
CREATE TABLE IF NOT EXISTS auto_close(
	"guild_id" int8 NOT NULL,
	"enabled" bool NOT NULL,
	"since_open_with_no_response" interval,
	"since_last_message" interval,
	PRIMARY KEY("guild_id")
);
`
}

func (a *AutoCloseTable) Get(guildId uint64) (settings AutoCloseSettings, e error) {
	query := `SELECT "enabled", "since_open_with_no_response", "since_last_message" FROM auto_close WHERE "guild_id" = $1;`
	if err := a.QueryRow(context.Background(), query, guildId).Scan(&settings.Enabled, &settings.SinceOpenWithNoResponse, &settings.SinceLastMessage); err != nil && err != pgx.ErrNoRows { // defaults to nil if no rows
		e = err
	}

	return
}

func (a *AutoCloseTable) Set(guildId uint64, settings AutoCloseSettings) (err error) {
	query := `INSERT INTO auto_close("guild_id", "enabled", "since_open_with_no_response", "since_last_message") VALUES($1, $2, $3, $4) ON CONFLICT("guild_id") DO UPDATE SET "enabled" = $2, "since_open_with_no_response" = $3, "since_last_message" = $4;`
	_, err = a.Exec(context.Background(), query, guildId, settings.Enabled, settings.SinceOpenWithNoResponse, settings.SinceLastMessage)
	return
}
