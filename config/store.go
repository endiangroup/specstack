package config

var onStoreValidations = []Validation{}

func Store(storer Storer, c *Config) error {
	isValid, err := IsValid(c, onStoreValidations...)
	if !isValid {
		return err
	}

	return storer.StoreConfig(c)
}
