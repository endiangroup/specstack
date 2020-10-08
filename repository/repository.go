package repository

import "io"

// Repository represents a version control repo
type Repository interface {
	Initialiser
	Configurer
	MetadataSyncer
	ObjectHasher
}

// Initialiser initialises a repo
type Initialiser interface {
	IsInitialised() bool
}

// Configurer allows for low level config settings
type Configurer interface {
	GetConfig(string) (string, error)
	SetConfig(string, string) error
	UnsetConfig(string) error
	AllConfig() (map[string]string, error)
}

// MetadataSyncer allows for low level metadata management
type MetadataSyncer interface {
	PrepareMetadataSync() error
	PullMetadata(from string) error
	PushMetadata(to string) error
}

type ObjectHasher interface {
	ObjectHash(io.Reader) (string, error)
	ObjectString(hash string) (string, error)
}
