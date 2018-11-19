package repository

import "github.com/endiangroup/specstack/persistence"

// Repository represents a version control repo
type Repository interface {
	Initialiser
	persistence.ConfigStorer
}

// Initialiser initialises a repo
type Initialiser interface {
	Init() error
	IsInitialised() bool
}
