package personas

import (
	"context"

	"github.com/endiangroup/specstack/config"
)

type Developer interface {
	ListConfiguration(context.Context) (map[string]string, error)
	GetConfiguration(context.Context, string) (string, error)
}

func NewDeveloper() *developer {
	return &developer{}
}

type developer struct {
}

func (d *developer) ListConfiguration(ctx context.Context) (map[string]string, error) {
	return config.FromContext(ctx).ToMap(), nil
}

func (d *developer) GetConfiguration(ctx context.Context, name string) (string, error) {
	return config.FromContext(ctx).Get(name)
}
