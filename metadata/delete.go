package metadata

import (
	"fmt"
	"io"
	"time"

	uuid "github.com/satori/go.uuid"
)

func Delete(storer Storer, key io.Reader, id uuid.UUID) error {
	var entries []*Entry
	if err := storer.ReadAllMetadata(key, &entries); err != nil {
		return err
	}

	var candidate *Entry
	for _, entry := range entries {
		if entry.Id == id {
			candidate = entry
		}
	}

	if candidate == nil {
		return fmt.Errorf("No entry for id %s", id)
	}

	candidate.DeletedAt = time.Now()

	return storer.StoreMetadata(key, candidate)
}
