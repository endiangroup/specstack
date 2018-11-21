package persistence

import (
	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/errors"
)

var (
	ErrNoConfigFound = errors.New("no config found")
)

func (store *Store) StoreConfig(c *config.Config) (*config.Config, error) {
	configMap := config.ToMap(c)

	errs := errors.Errors{}
	for key, value := range configMap {
		if err := store.ConfigStorer.SetConfig(key, value); err != nil {
			errs = errs.Append(err)
		}
	}

	if errs.Any() {
		return nil, errs
	}

	return c, nil
}

func (store *Store) LoadConfig() (*config.Config, error) {
	configMap, err := store.ConfigStorer.AllConfig()
	if err != nil {
		return nil, err
	}

	return config.NewFromMap(configMap)
}
