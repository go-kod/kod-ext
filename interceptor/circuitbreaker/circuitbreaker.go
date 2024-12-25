package circuitbreaker

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kod/kod-ext/internal/kerror"
	"github.com/go-kod/kod/interceptor"
	"github.com/sony/gobreaker"
)

type Config struct {
	Enable              bool          `json:"enable"`
	FailureRatio        float32       `json:"failure_ratio" default:"0.6"`
	MaxRequests         uint32        `json:"max_requests" default:"3"`
	Interval            time.Duration `json:"interval" default:"5s"`
	Timeout             time.Duration `json:"timeout" default:"10s"`
	ConsecutiveFailures uint32        `json:"consecutive_failures" default:"3"`
}

func (c Config) Interceptor() interceptor.Interceptor {
	if !c.Enable {
		return nil
	}

	breaker := gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
		MaxRequests: c.MaxRequests,
		Interval:    c.Interval,
		Timeout:     c.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float32(counts.TotalFailures) / float32(counts.Requests)
			return counts.ConsecutiveFailures >= c.ConsecutiveFailures && failureRatio >= c.FailureRatio
		},
	})

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		done, err := breaker.Allow()
		if err != nil {
			return fmt.Errorf("circuit breaker open: %w", err)
		}

		err = invoker(ctx, info, req, reply)

		done(!kerror.IsCritical(err))

		return err
	}
}
