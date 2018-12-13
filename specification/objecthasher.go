package specification

import "io"

type ObjectHasher interface {
	ObjectHash(io.Reader) (string, error)
}
