package specification

//go:generate stringer -type=SourceType
type SourceType int

const (
	SourceTypeFile SourceType = iota
	SourceTypeText
)

type Source struct {
	Type SourceType
	Body string
}

type Sourcer interface {
	Source() Source
}
