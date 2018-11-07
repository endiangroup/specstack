package persistence

import "github.com/endiangroup/specstack/repository"

func NewRepositoryStore(kvStore repository.ConfigStorer) *RepositoryStore {
	return &RepositoryStore{
		KVStore: kvStore,
	}
}

type RepositoryStore struct {
	KVStore repository.ConfigStorer
}
