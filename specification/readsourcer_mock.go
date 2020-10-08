package specification

import (
	"bytes"
	io "io"
)

type MockReadSourcer struct {
}

func (m *MockReadSourcer) ReadSource(a Sourcer) (io.Reader, error) {
	return bytes.NewBufferString(a.Source().Body), nil
}
