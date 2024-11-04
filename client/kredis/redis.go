package kredis

import (
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
)

type Client = redis.Client

type Config struct {
	Addr         string        `default:"-"`
	DialTimeout  time.Duration `default:"3s"`
	ReadTimeout  time.Duration `default:"1s"`
	WriteTimeout time.Duration `default:"1s"`
	Password     string        `default:""`
	DB           int           `default:"0"`
}

func (c Config) Build() *Client {
	if c.Addr == "" {
		panic("redis addr is required")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         c.Addr,
		Password:     c.Password,
		DB:           c.DB,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	return rdb
}
