package metadata

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type Entry struct {
	Id      uuid.UUID
	Created time.Time
	Deleted bool
	Name    string
	Value   string
}
