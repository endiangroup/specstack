package config

import (
	"strings"

	"github.com/endiangroup/specstack/errors"
)

func NewWithDefaults() *Config {
	return &Config{
		Project: NewProjectWithDefaults(),
	}
}

func New() *Config {
	return &Config{Project: newProject()}
}

func NewFromMap(configMap map[string]string) *Config {
	c := New()
	for key, value := range configMap {
		c.Set(key, value)
	}

	return c
}

type Config struct {
	Project *Project
}

func (c *Config) IsValid(validations ...Validation) (bool, error) {
	errs := errors.ValidationErrors{}

	for _, validation := range validations {
		if err := validation(c); err != nil {
			switch err := err.(type) {
			case *errors.ValidationField:
				errs = errs.Append(err)
			default:
				return false, err
			}
		}
	}

	if errs.Any() {
		return false, errs
	}

	return true, nil
}

func (c *Config) ToMap() map[string]string {
	return c.Project.ToMap("project.", map[string]string{})
}

func (c *Config) Set(name, value string) {
	nameParts := strings.Split(name, ".")

	switch nameParts[0] {
	case "project":
		c.Project.Set(strings.Join(nameParts[1:], "."), value)
	}
}
