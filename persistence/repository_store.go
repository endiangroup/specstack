package persistence

import "github.com/endiangroup/specstack/repository"

func NewRepositoryStore(kvStore repository.ConfigStorer) *RepositoryStore {
	return &RepositoryStore{
		ConfigStorer: kvStore,
	}
}

type RepositoryStore struct {
	ConfigStorer repository.ConfigStorer
}
