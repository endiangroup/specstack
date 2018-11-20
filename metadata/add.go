package metadata

import (
	"io"
	"time"

	uuid "github.com/satori/go.uuid"
)

func assertHeaders(entry *Entry) error {
	zeroId := uuid.UUID{}
	if entry.Id == zeroId {
		uid := uuid.NewV4()
		entry.Id = uid
	}

	zeroTime := time.Time{}
	if entry.CreatedAt == zeroTime {
		entry.CreatedAt = time.Now()
	}

	return nil
}

func Add(storer Storer, key io.Reader, entry *Entry) error {
	if err := assertHeaders(entry); err != nil {
		return err
	}
	return storer.StoreMetadata(key, entry)
}
