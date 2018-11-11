package specification

// A Reader represents the input for a specification. The read method
// returns a Specification, zero or more warnings, and a fatal error.
type Reader interface {
	Read() (*Specification, []error, error)
}
