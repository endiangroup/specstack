package persistence

import "github.com/endiangroup/specstack/repository"

func NewRepositoryStore(kvStore repository.KeyValueStorer) *RepositoryStore {
	return &RepositoryStore{
		KVStore: kvStore,
	}
}

type RepositoryStore struct {
	KVStore repository.KeyValueStorer
}
