package repository

type Writer interface {
	Init() error
}

type Reader interface {
	IsInitialised() bool
	ConfigGetRegex(string) (string, error)
}

type ReadWriter interface {
	Reader
	Writer
}
