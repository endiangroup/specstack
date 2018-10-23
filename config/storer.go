package config

type Storer interface {
	CreateConfig(*Config) (*Config, error)
	LoadConfig() (*Config, error)
	StoreConfig(*Config) error
}
