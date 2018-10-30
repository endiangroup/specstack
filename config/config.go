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

func NewFromMap(configMap map[string]string) *Config {
	c := New()
	for key, value := range configMap {
		Set(c, key, value)
	}

	return c
}

type Config struct {
	Project *Project
	User    *User
}
