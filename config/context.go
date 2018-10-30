package config

import "context"

const (
	contextConfigKey = "config"
)

func InContext(ctx context.Context, c *Config) context.Context {
	return context.WithValue(ctx, contextConfigKey, c)
}

func FromContext(ctx context.Context) *Config {
	if c, ok := ctx.Value(contextConfigKey).(*Config); ok && c != nil {
		return c
	}

	return Empty
}
