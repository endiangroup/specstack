package metadata

import "io"

type GetterSetter interface {
	GetMetadata(key io.Reader, output interface{}) error
	SetMetadata(key io.Reader, value interface{}) error
}
