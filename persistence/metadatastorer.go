package persistence

import "io"

type MetadataStorer interface {
	GetMetadata(key io.Reader) ([][]byte, error)
	SetMetadata(key io.Reader, value []byte) error
}
