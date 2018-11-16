package personas

import (
	"context"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/repository"
)

type Developer interface {
	ListConfiguration(context.Context) (map[string]string, error)
	GetConfiguration(context.Context, string) (string, error)
	SetConfiguration(context.Context, string, string) error
}

func NewDeveloper(
	configStore config.Storer,
	repo repository.Repository,
) *developer {
	return &developer{
		configStore: configStore,
		repo:        repo,
	}
}

type developer struct {
	configStore config.Storer
	repo        repository.Repository
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

	_, err = config.Store(d.configStore, c)
	return err
}
