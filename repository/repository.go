package repository

// Repository represents a version control repo
type Repository interface {
	Initialiser
	ConfigStorer
}

// Initialiser initialises a repo
type Initialiser interface {
	Init() error
	IsInitialised() bool
}

// ConfigStorer is used to read and write key-value config
type ConfigStorer interface {
	GetConfig(string) (string, error)
	SetConfig(string, string) error
	UnsetConfig(string) error
	AllConfig() (map[string]string, error)
}
