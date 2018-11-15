package repository

import "io"

// Repository represents a version control repo
type Repository interface {
	Initialiser
	MetadataStorer
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

// MetadataStorer is used to read and write repo metadata.
// In Git, this takes the form of git notes.
// GetMetaData requires an underlying slice for its output
// argument, otherwise it cannot process the data.
type MetadataStorer interface {
	GetMetadata(key io.Reader, output interface{}) error
	SetMetadata(key io.Reader, value interface{}) error
}
