package faultinject

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-kod/kod/interceptor"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type FaultInjectConfig struct {
	Enable     bool          `json:"enable"`
	Delay      time.Duration `json:"delay"`
	Error      string        `json:"error"`
	Percentage uint32        `json:"percentage" default:"100"`
}

func (c FaultInjectConfig) Interceptor() interceptor.Interceptor {
	if !c.Enable || c.Percentage == 0 {
		return nil
	}

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		// calculate the percentage of requests to inject a fault into
		if c.Percentage < 100 {
			if rand.Intn(100) > int(c.Percentage) {
				return invoker(ctx, info, req, reply)
			}
		}

		if c.Delay > 0 {
			span := trace.SpanFromContext(ctx)
			if span.SpanContext().IsValid() {
				span.AddEvent("fault inject delay",
					trace.WithAttributes(attribute.String("delay", c.Delay.String())))
			}

			time.Sleep(c.Delay)
		}

		if c.Error != "" {
			return fmt.Errorf("fault inject: %s", c.Error)
		}

		return invoker(ctx, info, req, reply)
	}
}
