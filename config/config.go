package config

type ConfigReader interface {
	List() (string, error)
}
