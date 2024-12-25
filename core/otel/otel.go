package otel

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

type Config struct {
	EnableTrace  bool
	EnableMetric bool
	EnableLog    bool
}

func (c Config) Init(ctx context.Context, k *kod.Kod) error {
	resource := lo.Must(sdkresource.New(ctx,
		sdkresource.WithFromEnv(),
		sdkresource.WithTelemetrySDK(),
		sdkresource.WithHost(),
		sdkresource.WithContainer(),
		sdkresource.WithAttributes(
			semconv.ServiceNameKey.String(k.Config().Name),
			semconv.ServiceVersionKey.String(k.Config().Version),
			semconv.DeploymentEnvironmentKey.String(k.Config().Env),
		)),
	)

	// configure trace provider
	if c.EnableTrace {
		spanExporter := lo.Must(autoexport.NewSpanExporter(ctx))
		spanProvider := sdktrace.NewTracerProvider(
			sdktrace.WithBatcher(spanExporter),
			sdktrace.WithResource(resource),
		)

		otel.SetTracerProvider(spanProvider)

		propagator := propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{}, propagation.Baggage{},
		)

		otel.SetTextMapPropagator(propagator)
		k.Defer("OTEL-Trace", spanProvider.Shutdown)
	}

	// configure metric provider
	if c.EnableMetric {
		lo.Must0(host.Start())
		lo.Must0(runtime.Start())

		metricReader := lo.Must(autoexport.NewMetricReader(ctx))
		metricProvider := sdkmetric.NewMeterProvider(
			sdkmetric.WithReader(metricReader),
			sdkmetric.WithResource(resource),
		)

		otel.SetMeterProvider(metricProvider)
		k.Defer("OTEL-Metric", metricProvider.Shutdown)
	}

	// configure log provider
	if c.EnableLog {
		logExporter := lo.Must(autoexport.NewLogExporter(ctx))
		loggerProvider := sdklog.NewLoggerProvider(
			sdklog.WithProcessor(
				sdklog.NewBatchProcessor(logExporter),
			),
			sdklog.WithResource(resource),
		)

		global.SetLoggerProvider(loggerProvider)
		k.Defer("OTEL-Log", loggerProvider.Shutdown)
	}

	return nil
}
