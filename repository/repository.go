package repository

type Repository interface {
	Initialiser
	ConfigStorer
}

type Initialiser interface {
	Init() error
	IsInitialised() bool
}

type ConfigStorer interface {
	GetConfig(string) (string, error)
	SetConfig(string, string) error
	UnsetConfig(string) error
	AllConfig() (map[string]string, error)
}
