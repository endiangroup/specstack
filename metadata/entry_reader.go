package metadata

import "time"

// ID represents the address or other ID of the metadata in the repo
type ID string

func (i ID) String() string {
	return string(i)
}

// An EntryReader reads a metadata enty
type EntryReader interface {
	Id() ID
	Created() time.Time
	Read() (interface{}, error)
}
