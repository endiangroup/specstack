package personas

import (
	"context"

	"github.com/endiangroup/specstack/config"
)

type Developer interface {
	ListConfiguration(context.Context) (map[string]string, error)
	GetConfiguration(context.Context, string) (string, error)
	SetConfiguration(context.Context, string, string) error
}

func NewDeveloper(configStore config.Storer) *developer {
	return &developer{
		configStore: configStore,
	}
}

type developer struct {
	configStore config.Storer
}

func (d *developer) ListConfiguration(ctx context.Context) (map[string]string, error) {
	return config.ToMap(config.FromContext(ctx)), nil
}

func (d *developer) GetConfiguration(ctx context.Context, name string) (string, error) {
	return config.Get(config.FromContext(ctx), name)
}

func (d *developer) SetConfiguration(ctx context.Context, name, value string) error {
	c := config.FromContext(ctx)

	err := config.Set(c, name, value)
	if err != nil {
		return err
	}

	return config.Store(d.configStore, c)
}
