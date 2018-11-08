package metadata

// An EntryWriter can update a metadata entry
type EntryWriter interface {
	Write(interface{}) error
}
