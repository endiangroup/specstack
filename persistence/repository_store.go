package persistence

import "github.com/endiangroup/specstack/repository"

func NewRepositoryStore(repo repository.ReadWriter) *RepositoryStore {
	return &RepositoryStore{
		Repo: repo,
	}
}

type RepositoryStore struct {
	Repo repository.ReadWriter
}
