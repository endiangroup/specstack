package metadata

import (
	"io"
)

type Storer interface {
	StoreMetadata(key io.Reader, value interface{}) error
	ReadAllMetadata(key io.Reader, into interface{}) error
}
