package persistence

type ConfigStorer interface {
	GetConfig(string) (string, error)
	SetConfig(string, string) error
	UnsetConfig(string) error
	AllConfig() (map[string]string, error)
}
