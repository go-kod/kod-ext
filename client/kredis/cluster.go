package kredis

import (
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
)

type ClusterClient = redis.ClusterClient

type ClusterConfig struct {
	Addrs        []string      `default:"-"`
	DialTimeout  time.Duration `default:"3s"`
	ReadTimeout  time.Duration `default:"1s"`
	WriteTimeout time.Duration `default:"1s"`
	Password     string        `default:""`
	DB           int           `default:"0"`
}

func (c ClusterConfig) Build() *ClusterClient {
	if len(c.Addrs) == 0 {
		panic("redis cluster addrs is required")
	}

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        c.Addrs,
		Password:     c.Password,
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
