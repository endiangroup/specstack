package metadata

import (
	"io"

	uuid "github.com/satori/go.uuid"
)

type Storer interface {
	Store(io.Reader, *Entry) error
	Delete(io.Reader, uuid.UUID) error
}

type Reader interface {
	Read(io.Reader) ([]*Entry, error)
}

type ReadStorer interface {
	Storer
	Reader
}
