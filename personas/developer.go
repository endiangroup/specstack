package personas

import (
	"context"
	"io"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/metadata"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/specification"
)

type Developer interface {
	ListConfiguration(context.Context) (map[string]string, error)
	GetConfiguration(context.Context, string) (string, error)
	SetConfiguration(context.Context, string, string) error
	AddMetadataToStory(context.Context, *specification.Story, io.Reader, string, string) error
	AddMetadataToScenario(context.Context, *specification.Scenario, io.Reader, string, string) error
}

func NewDeveloper(
	store *persistence.Store,
) *developer {
	return &developer{
		store: store,
	}
}

type developer struct {
	store *persistence.Store
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

	_, err = config.Store(d.store, c)
	return err
}

func (d *developer) AddMetadataToStory(
	ctx context.Context,
	story *specification.Story,
	object io.Reader,
	name, value string,
) error {
	return metadata.Add(d.store, object, metadata.NewKeyValue(name, value))
}

func (d *developer) AddMetadataToScenario(
	ctx context.Context,
	story *specification.Scenario,
	object io.Reader,
	name, value string,
) error {
	return metadata.Add(d.store, object, metadata.NewKeyValue(name, value))
}
