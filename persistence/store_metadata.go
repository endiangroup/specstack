package persistence

import (
	"fmt"
	io "io"
	"sort"
	"time"

	"github.com/endiangroup/specstack/metadata"
	uuid "github.com/satori/go.uuid"
)

func (r *Store) assertHeaders(entry *metadata.Entry) error {
	zeroId := uuid.UUID{}
	if entry.Id == zeroId {
		uid := uuid.NewV4()
		entry.Id = uid
	}

	zeroTime := time.Time{}
	if entry.Created == zeroTime {
		entry.Created = time.Now()
	}

	return nil
}

func (r *Store) StoreMetadata(key io.Reader, entry *metadata.Entry) error {
	if err := r.assertHeaders(entry); err != nil {
		return err
	}
	return r.MetadataStorer.SetMetadata(key, entry)
}

func (r *Store) DeleteMetadata(key io.Reader, id uuid.UUID) error {
	var entries []*metadata.Entry
	if err := r.MetadataStorer.GetMetadata(key, &entries); err != nil {
		return err
	}

	var candidate *metadata.Entry
	for _, entry := range entries {
		if entry.Id == id {
			candidate = entry
		}
	}

	if candidate == nil {
		return fmt.Errorf("No entry for id %s", id)
	}

	candidate.Deleted = true

	return r.MetadataStorer.SetMetadata(key, candidate)
}

func (r *Store) ReadMetadata(key io.Reader) ([]*metadata.Entry, error) {

	// TODO! Merge on constraints. The current implemenation
	// cares about unique names only.
	var outputs []*metadata.Entry

	if err := r.MetadataStorer.GetMetadata(key, &outputs); err != nil {
		return nil, err
	}

	entryMap := make(map[string]*metadata.Entry)

	// Outputs are returned in chronological order,
	// so we can step through them an take the most
	// recent as canon.
	for _, output := range outputs {
		if !output.Deleted {
			entryMap[output.Name] = output
		}
	}

	//nolint:prealloc
	var final []*metadata.Entry
	for _, e := range entryMap {
		final = append(final, e)
	}

	sort.Slice(final, func(i, j int) bool {
		return final[i].Name < final[j].Name
	})

	return final, nil
}
