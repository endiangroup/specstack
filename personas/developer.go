package personas

import "github.com/endiangroup/specstack/config"

type Developer interface {
	ListConfiguration() (map[string]string, error)
}

func NewDeveloper(configStore config.Storer) *developer {
	return &developer{
		configStore: configStore,
	}
}

type developer struct {
	configStore config.Storer
}

func (d *developer) ListConfiguration() (map[string]string, error) {
	config, err := d.configStore.LoadConfig()
	if err != nil {
		return nil, err
	}

	return config.ToMap(), nil
}
