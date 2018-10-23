package config

var onStoreValidations = []Validation{}

func Store(storer Storer, c *Config) (*Config, error) {
	isValid, err := IsValid(c, onStoreValidations...)
	if !isValid {
		return nil, err
	}

	return storer.StoreConfig(c)
}
