package ksql

import (
	"database/sql"
	"time"

	"github.com/samber/lo"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
)

type Config struct {
	DriverName      string        `default:"sqlite3"`
	Dsn             string        `default:"-"`
	MaxIdleConns    int           `default:"25"`
	MaxOpenConns    int           `default:"500"`
	ConnMaxLifetime time.Duration `default:"1h"`
	ConnMaxIdleTime time.Duration `default:"1h"`
}

func (c Config) Build() *sql.DB {
	if c.Dsn == "" {
		panic("sql: DSN is empty")
	}

	db := lo.Must(otelsql.Open(c.DriverName, c.Dsn,
		otelsql.WithDBSystem(c.DriverName),
	))

	db.SetMaxIdleConns(c.MaxIdleConns)
	db.SetMaxOpenConns(c.MaxOpenConns)
	db.SetConnMaxLifetime(c.ConnMaxLifetime)
	db.SetConnMaxIdleTime(c.ConnMaxIdleTime)

	return db
}
