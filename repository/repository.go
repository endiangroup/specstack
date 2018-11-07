package repository

// Repository represents a version control repo
type Repository interface {
	Initialiser
	KeyValueStorer
	MetadataStorer
}

// Initialiser initialises a repo
type Initialiser interface {
	Init() error
	IsInitialised() bool
}

// KeyValueStorer is used to read and write key-value config
type KeyValueStorer interface {
	Get(string) (string, error)
	Set(string, string) error
	Unset(string) error
	All() (map[string]string, error)
}

// MetadataStorer is used to read and write repo metadata.
// In Git, this takes the form of git notes.
type MetadataStorer interface {
	GetMetadata(key string) (string, error)
	SetMetadata(key, value string) error
}
