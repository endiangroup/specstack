package metadata

import (
	"io"
	"sort"
)

func ReadAll(storer Storer, key io.Reader) ([]*Entry, error) {

	// TODO! Merge on constraints. The current implemenation
	// cares about unique names only.
	var outputs []*Entry

	if err := storer.ReadAllMetadata(key, &outputs); err != nil {
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
