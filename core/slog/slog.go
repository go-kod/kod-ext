package slog

import (
	"context"
	"log/slog"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/internal/kslog"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Config struct {
	LogLevel slog.Level `default:"info"`
}

func (c Config) Init(ctx context.Context) error {
	k := kod.FromContext(ctx)

	handler := kslog.NewLevelHandler(c.LogLevel)(
		otelslog.NewHandler(k.Config().Name),
	)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return nil
}
