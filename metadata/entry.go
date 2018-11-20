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

func (e *Entry) IsDeleted() bool {
	return e.DeletedAt.IsZero()
}
