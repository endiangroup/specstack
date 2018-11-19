package metadata

import (
	"io"

	uuid "github.com/satori/go.uuid"
)

type Storer interface {
	ReadMetadata(key io.Reader) ([]*Entry, error)
	StoreMetadata(key io.Reader, entry *Entry) error
	DeleteMetadata(key io.Reader, id uuid.UUID) error
}
