package config

import "github.com/endiangroup/specstack/repository"

const (
	namespace      = "specstack"
	namespaceRegex = "^" + namespace + `\.`
)

func NewRepositoryConfig(repo repository.ReadWriter) RepositoryConfig {
	return RepositoryConfig{
		Repo: repo,
	}
}

type RepositoryConfig struct {
	Repo repository.ReadWriter
}

func (r RepositoryConfig) List() (string, error) {
	return r.Repo.ConfigGetRegex(namespaceRegex)
}
