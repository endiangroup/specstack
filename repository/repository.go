package repository

type Repository interface {
	Initialiser
	KeyValueStorer
}

type Initialiser interface {
	Init() error
	IsInitialised() bool
}

type KeyValueStorer interface {
	Get(string) (string, error)
	Set(string, string) error
	Unset(string) error
	All() (map[string]string, error)
}
