package repository

type Writer interface {
	Init() error
	ConfigSet(string, string) error
	ConfigUnset(string) error
}

type Reader interface {
	IsInitialised() bool
	ConfigGetAll() (map[string]string, error)
	ConfigGet(string) (string, error)
}

type ReadWriter interface {
	Reader
	Writer
}
