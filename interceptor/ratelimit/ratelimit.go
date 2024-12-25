package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kod/kod/interceptor"
	"github.com/juju/ratelimit"
)

type Config struct {
	Enable  bool          `json:"enable"`
	QPS     int64         `json:"qps"`
	MaxWait time.Duration `json:"max_wait" default:"5s"`
}

func (c Config) Build() interceptor.Interceptor {
	if !c.Enable || c.QPS <= 0 {
		return nil
	}

	limiter := ratelimit.NewBucketWithRate(float64(c.QPS), c.QPS)

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		wait := limiter.WaitMaxDuration(1, c.MaxWait)
		if !wait {
			return fmt.Errorf("rate limit exceeded")
		}

		return invoker(ctx, info, req, reply)
	}
}
