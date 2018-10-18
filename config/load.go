package config

var onLoadValidations = []Validation{}

func Load(storer Storer) (*Config, error) {
	config, err := storer.LoadConfig()
	if err != nil {
		return nil, err
	}

	isValid, err := config.IsValid(onLoadValidations...)
	if !isValid {
		return nil, err
	}

	return config, nil
}
