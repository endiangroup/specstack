package metadata

// Writer allows the creation (via appending) and deletion of metadata
type Writer interface {
	Append() (EntryReadWriter, error)
	Delete(EntryReader) error
}
