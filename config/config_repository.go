package config

import "github.com/endiangroup/specstack/repository"

const (
	namespace      = "specstack"
	namespaceRegex = "^" + namespace + `\.`
)

func NewRepositoryConfig(repo repository.RepositoryReadWriter) RepositoryConfig {
	return RepositoryConfig{
		Repo: repo,
	}
}

type RepositoryConfig struct {
	Repo repository.RepositoryReadWriter
}

func (r RepositoryConfig) List() (string, error) {
	return r.Repo.ConfigRegex(namespaceRegex)
}
