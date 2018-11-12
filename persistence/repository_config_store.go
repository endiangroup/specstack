package persistence

import (
	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/errors"
	"github.com/endiangroup/specstack/repository"
)

var (
	ErrNoConfigFound = errors.New("no config found")
)

func (store *RepositoryStore) StoreConfig(c *config.Config) (*config.Config, error) {
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

func (store *RepositoryStore) LoadConfig() (*config.Config, error) {
	configMap, err := store.ConfigStorer.AllConfig()
	if err != nil {
		if _, ok := err.(repository.GitConfigMissingKeyErr); ok {
			return nil, ErrNoConfigFound
		}

		return nil, err
	}

	return config.NewFromMap(configMap)
}
