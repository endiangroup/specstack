package config

var onCreateValidations = []Validation{
	UserNameCannotBeBlank,
	UserEmailCannotBeBlank,
}

func CreateDefault(storer Storer) (*Config, error) {
	return Create(storer, NewWithDefaults())
}

func Create(storer Storer, c *Config) (*Config, error) {
	isValid, err := IsValid(c, onCreateValidations...)
	if !isValid {
		return nil, err
	}

	return storer.StoreConfig(c)
}
