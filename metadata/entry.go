package metadata

import (
	"time"
)

type Entry struct {
	CreatedAt time.Time
	Name      string
	Value     string
}

func New() *Entry {
	return &Entry{}
}

func NewKeyValue(key, value string) *Entry {
	e := New()
	e.Name = key
	e.Value = value
	return e
}
