package config

var (
	Empty = New()
)

func NewWithDefaults() *Config {
	return &Config{
		Project: newProjectWithDefaults(),
		User:    newUserWithDefaults(),
	}
}

func New() *Config {
	return &Config{
		Project: newProject(),
		User:    newUser(),
	}
}

func NewFromMap(configMap map[string]string) (*Config, error) {
	c := New()
	for key, value := range configMap {
		if err := Set(c, key, value); err != nil {
			return nil, err
		}
	}

	return c, nil
}

type Config struct {
	Project *Project
	User    *User
}
