package specification

import (
	"crypto/md5"
	"fmt"
	io "io"
)

type MockObjectHasher struct {
}

func (m *MockObjectHasher) ObjectHash(r io.Reader) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, r)
	return fmt.Sprintf("%x", h.Sum(nil)), err
}
