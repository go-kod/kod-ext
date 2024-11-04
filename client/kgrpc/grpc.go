package kgrpc

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/go-kod/kod-ext/registry"
)

type ClientConn = grpc.ClientConn

type Config struct {
	Target  string        `default:"-"`
	Timeout time.Duration `default:"3s"`

	registry registry.Registry
}

func (c Config) WithRegistry(r registry.Registry) Config {
	c.registry = r
	return c
}

func (c Config) Build(opts ...grpc.DialOption) *ClientConn {
	if c.Target == "" {
		panic("grpc target is required")
	}

	ctx := context.Background()

	defaultOpts := []grpc.DialOption{
		grpc.WithNoProxy(),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	if c.registry != nil {
		builder, err := c.registry.ResolveBuilder(ctx)
		if err != nil {
			panic(err)
		}
		defaultOpts = append(defaultOpts, grpc.WithResolvers(builder))
	}

	defaultOpts = append(defaultOpts, opts...)

	cc, err := grpc.NewClient(c.Target, defaultOpts...)
	if err != nil {
		panic(err)
	}

	return cc
}
