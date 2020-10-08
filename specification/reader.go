package specification

import (
	"io"

	"github.com/endiangroup/specstack/errors"
)

type ReadSourcer interface {
	ReadSource(Sourcer) (io.Reader, error)
}

// A Reader represents the input for a specification. The read method
// returns a Specification, zero or more warnings, and a fatal error.
type Reader interface {
	ReadSourcer
	Read() (*Specification, errors.Warnings, error)
}
