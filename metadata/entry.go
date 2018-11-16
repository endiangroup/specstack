package metadata

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

//go:generate stringer -type=Status
type Status int

const (
	StatusNormal  Status = 0
	StatusDeleted        = iota
)

type Entry struct {
	Id      uuid.UUID
	Created time.Time
	Status  Status
	Name    string
	Value   string
}
