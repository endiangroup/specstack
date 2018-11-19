package persistence

import "io"

type MetadataStorer interface {
	GetMetadata(key io.Reader, output interface{}) error
	SetMetadata(key io.Reader, value interface{}) error
}
