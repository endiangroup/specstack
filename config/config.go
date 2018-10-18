package config

var (
	EmptyConfig = New()
)

func NewWithDefaults() *Config {
	return &Config{
		Project: newProjectWithDefaults(),
	}
}

func New() *Config {
	return &Config{Project: newProject()}
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
}
