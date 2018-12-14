package specification

import (
	"crypto/md5"
	"fmt"
	io "io"
	"io/ioutil"
)

type MockObjectHasher struct {
}

func (m *MockObjectHasher) ObjectHash(r io.Reader) (string, error) {
	h := md5.New()
	_, err := io.Copy(h, r)
	return fmt.Sprintf("%x", h.Sum(nil)), err
}

type MockObjectHasherPlaintext struct {
}

func (m *MockObjectHasherPlaintext) ObjectHash(r io.Reader) (string, error) {
	data, err := ioutil.ReadAll(r)
	return string(data), err
}
