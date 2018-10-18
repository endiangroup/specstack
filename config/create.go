package config

var onCreateValidations = []Validation{}

func CreateDefault(storer Storer) (*Config, error) {
	return Create(storer, NewWithDefaults())
}

func Create(storer Storer, c *Config) (*Config, error) {
	isValid, err := IsValid(c, onCreateValidations...)
	if !isValid {
		return nil, err
	}

	return storer.CreateConfig(c)
}
