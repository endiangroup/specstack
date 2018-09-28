package repository

type Writer interface {
	Init() error
}

type Reader interface {
	IsInitialised() bool
	ConfigRegex(string) (string, error)
}

type ReadWriter interface {
	Reader
	Writer
}
