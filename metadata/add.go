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

func Add(storer Storer, key io.Reader, entries ...*Entry) error {
	for _, entry := range entries {
		if err := assertHeaders(entry); err != nil {
			return err
		}
		if err := storer.StoreMetadata(key, entry); err != nil {
			return err
		}
	}
	return nil
}
