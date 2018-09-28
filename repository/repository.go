package repository

type RepositoryWriter interface {
	Init() error
}

type RepositoryReader interface {
	IsInitialised() bool
	ConfigRegex(string) (string, error)
}

type RepositoryReadWriter interface {
	RepositoryReader
	RepositoryWriter
}
