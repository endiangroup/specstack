package config

var onCreateValidations = []Validation{}

func CreateDefault(storer Storer) (*Config, error) {
	return Create(storer, NewWithDefaults())
}

func Create(storer Storer, config *Config) (*Config, error) {
	isValid, err := config.IsValid(onCreateValidations...)
	if !isValid {
		return nil, err
	}

	return storer.CreateConfig(config)
}
