package metadata

// Storer allows the creation (via appending) and deletion of metadata
type Storer interface {
	Append() (EntryReadWriter, error)
	Delete(EntryReader) error
}
