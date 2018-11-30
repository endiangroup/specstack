package repository

// Repository represents a version control repo
type Repository interface {
	Initialiser
	Configurer
}

// Initialiser initialises a repo
type Initialiser interface {
	Init() error
	IsInitialised() bool
}

type Configurer interface {
	GetConfig(string) (string, error)
	SetConfig(string, string) error
	UnsetConfig(string) error
	AllConfig() (map[string]string, error)
}
