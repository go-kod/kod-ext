package metric

import (
	"context"
	"time"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
)

var (
	methodCounts = lo.Must(otel.Meter(kod.PkgPath).Int64Counter("kod.component.count",
		metric.WithDescription("Count of Kod component method invocations"),
	))

	methodErrors = lo.Must(otel.Meter(kod.PkgPath).Int64Counter("kod.component.error",
		metric.WithDescription("Count of Kod component method invocations that result in an error"),
	))

	methodDurations = lo.Must(otel.Meter(kod.PkgPath).Float64Histogram("kod.component.duration",
		metric.WithDescription("Duration, in microseconds, of Kod component method execution"),
		metric.WithUnit("ms"),
	))
)

type Config struct{}

// Build returns an interceptor that adds OpenTelemetry metrics to the context.
func (c Config) Build() interceptor.Interceptor {
	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) (err error) {
		now := time.Now()

		err = invoker(ctx, info, req, reply)

		as := attribute.NewSet(
			attribute.String("method", info.FullMethod),
		)

		if err != nil {
			methodErrors.Add(ctx, 1, metric.WithAttributeSet(as))
		}

		methodCounts.Add(ctx, 1, metric.WithAttributeSet(as))
		methodDurations.Record(ctx, float64(time.Since(now).Milliseconds()), metric.WithAttributeSet(as))

		return err
	}
}
