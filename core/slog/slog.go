package slog

import (
	"context"
	"log/slog"

	"github.com/creasty/defaults"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod-ext/internal/kslog"
	"go.opentelemetry.io/contrib/bridges/otelslog"
)

type Config struct {
	LogLevel slog.Level `json:"log_level" default:"info"`
}

func (c Config) Init(ctx context.Context, k *kod.Kod) error {
	defaults.Set(&c)

	handler := kslog.NewLevelHandler(c.LogLevel)(
		otelslog.NewHandler(k.Config().Name),
	)
	logger := slog.New(handler)

	slog.SetDefault(logger)

	return nil
}
