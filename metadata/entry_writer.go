package metadata

// An EntryWriter can update a metadata entry
type EntryWriter interface {
	SetData(interface{}) error
}
