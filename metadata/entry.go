package metadata

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Entry struct {
	Id        uuid.UUID
	CreatedAt time.Time
	DeletedAt time.Time
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

func (e *Entry) IsDeleted() bool {
	return e.DeletedAt.IsZero()
}
