package metadata

import (
	"io"
	"time"
)

func assertHeaders(entry *Entry) error {
	if entry.CreatedAt.IsZero() {
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
