package config

type Reader interface {
	List() (string, error)
}
