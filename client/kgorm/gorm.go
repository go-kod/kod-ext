package kgorm

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

type Config struct {
	Dsn string `default:"-"`
}

func (c Config) Build() *DB {
	if c.Dsn == "" {
		panic("gorm dsn is required")
	}

	db, err := gorm.Open(mysql.Open(c.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}

	return db
}

type (
	// DB is an alias for gorm.DB.
	DB = gorm.DB
)
