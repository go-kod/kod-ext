package timeout

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kod/kod/interceptor"
)

type Config struct {
	Enable   bool          `json:"enable"`
	Duration time.Duration `json:"duration" default:"5s"`
}

func (c Config) Build() interceptor.Interceptor {
	if !c.Enable {
		return nil
	}

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		ctx, cancel := context.WithTimeoutCause(ctx, c.Duration, fmt.Errorf("timeout exceeded"))
		defer cancel()

		return invoker(ctx, info, req, reply)
	}
}
