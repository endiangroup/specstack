package config

type Storer interface {
	LoadConfig() (*Config, error)
	StoreConfig(*Config) (*Config, error)
}
