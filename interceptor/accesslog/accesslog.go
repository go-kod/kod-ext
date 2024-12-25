package accesslog

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-kod/kod-ext/internal/kerror"
	"github.com/go-kod/kod/interceptor"
)

type Config struct{}

// Interceptor returns a server interceptor that logs requests and responses.
func (c Config) Build() interceptor.Interceptor {
	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		now := time.Now()

		err := invoker(ctx, info, req, reply)

		attrs := []slog.Attr{
			slog.Any("req", req),
			slog.Any("reply", reply),
			slog.String("method", info.FullMethod),
			slog.Int64("cost", time.Since(now).Milliseconds()),
		}

		level := slog.LevelInfo
		if err != nil {
			level = slog.LevelError
			if kerror.IsBusiness(err) {
				level = slog.LevelWarn
			}

			attrs = append(attrs, slog.String("error", err.Error()))
		}

		// check if impl L(ctx context.Context) method
		if l, ok := info.Impl.(interface {
			L(context.Context) *slog.Logger
		}); ok {
			l.L(ctx).LogAttrs(ctx, level, "access_log", attrs...)
		}

		return err
	}
}
