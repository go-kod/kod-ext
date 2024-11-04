package http

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client = http.Client

type ClientConfig struct {
	Address string        `default:"-"`
	Timeout time.Duration `default:"3s"`
}

func (c ClientConfig) Build() *Client {
	if c.Address == "" {
		panic("http address is required")
	}

	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.Proxy = nil

	return &http.Client{
		Transport: otelhttp.NewTransport(defaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			})),
		Timeout: c.Timeout,
	}
}
