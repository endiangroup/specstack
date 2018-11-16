package metadata

import (
	"fmt"
	"io"
	"sort"
	"time"

	uuid "github.com/satori/go.uuid"
)

type readStorer struct {
	gs GetterSetter
}

func New(gs GetterSetter) ReadStorer {
	return &readStorer{
		gs: gs,
	}
}

func (r *readStorer) assertHeaders(entry *Entry) error {
	zeroId := uuid.UUID{}
	if entry.Id == zeroId {
		uid, err := uuid.NewV4()
		if err != nil {
			return err
		}
		entry.Id = uid
	}

	zeroTime := time.Time{}
	if entry.Created == zeroTime {
		entry.Created = time.Now()
	}

	return nil
}

func (r *readStorer) Store(key io.Reader, entry *Entry) error {
	if err := r.assertHeaders(entry); err != nil {
		return err
	}
	return r.gs.SetMetadata(key, entry)
}

func (r *readStorer) Delete(key io.Reader, id uuid.UUID) error {
	var entries []*Entry
	if err := r.gs.GetMetadata(key, &entries); err != nil {
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

	candidate.Status = StatusDeleted

	return r.gs.SetMetadata(key, candidate)
}

func (r *readStorer) Read(key io.Reader) ([]*Entry, error) {

	// TODO! Merge on constraints. The current implemenation
	// cares about unique names only.
	var outputs []*Entry

	if err := r.gs.GetMetadata(key, &outputs); err != nil {
		return nil, err
	}

	entryMap := make(map[string]*Entry)

	// Outputs are returned in chronological order,
	// so we can step through them an take the most
	// recent as canon.
	for _, output := range outputs {
		entryMap[output.Name] = output
	}

	//nolint:prealloc
	var final []*Entry
	for _, e := range entryMap {
		final = append(final, e)
	}

	sort.Slice(final, func(i, j int) bool {
		return final[i].Name < final[j].Name
	})

	return final, nil
}
