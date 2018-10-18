package config

var onLoadValidations = []Validation{}

func Load(storer Storer) (*Config, error) {
	c, err := storer.LoadConfig()
	if err != nil {
		return nil, err
	}

	isValid, err := IsValid(c, onLoadValidations...)
	if !isValid {
		return nil, err
	}

	return c, nil
}
