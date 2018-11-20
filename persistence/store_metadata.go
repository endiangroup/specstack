package persistence

import (
	"encoding/json"
	"fmt"
	io "io"
)

func (r *Store) StoreMetadata(key io.Reader, value interface{}) error {
	jsn, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to serialise metadata value: %s", err)
	}
	return r.MetadataStorer.SetMetadata(key, jsn)
}

func (r *Store) ReadAllMetadata(key io.Reader, into interface{}) error {
	encoded, err := r.MetadataStorer.GetMetadata(key)
	if err != nil {
		return fmt.Errorf("failed to get raw metadata: %s", err)
	}

	raw := []json.RawMessage{}
	for _, v := range encoded {
		raw = append(raw, json.RawMessage(v))
	}

	j, err := json.Marshal(raw)
	if err != nil {
		return fmt.Errorf("failed to marshal notes JSON: %s", err)
	}

	return json.Unmarshal(j, into)
}
