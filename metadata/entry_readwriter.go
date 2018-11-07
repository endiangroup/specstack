package metadata

// An EntryReadWriter can read or update a metadata entry
type EntryReadWriter interface {
	EntryReader
	EntryWriter
}
