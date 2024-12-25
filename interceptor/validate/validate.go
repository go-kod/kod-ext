package validate

import (
	"context"
	"fmt"

	"github.com/go-kod/kod/interceptor"
	"github.com/go-playground/validator/v10"
)

type Config struct{}

func (c Config) Build() interceptor.Interceptor {
	validate := validator.New()

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		for _, v := range req {
			if err := validate.Struct(v); err != nil {
				return fmt.Errorf("validate failed: %w", err)
			}
		}

		return invoker(ctx, info, req, reply)
	}
}
