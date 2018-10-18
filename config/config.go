package config

import (
	"context"
	"fmt"
	"strings"

	"github.com/endiangroup/specstack/errors"
)

const (
	contextConfigKey = "config"
)

var (
	EmptyConfig = New()
)

func InContext(ctx context.Context, c *Config) context.Context {
	return context.WithValue(ctx, contextConfigKey, c)
}

func FromContext(ctx context.Context) *Config {
	if c, ok := ctx.Value(contextConfigKey).(*Config); ok && c != nil {
		return c
	}

	return EmptyConfig
}

type ErrKeyNotFound string

func (err ErrKeyNotFound) Error() string {
	return fmt.Sprintf("no config key '%s' found", string(err))
}

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

func (c *Config) Get(name string) (string, error) {
	nameParts := strings.Split(name, ".")

	switch nameParts[0] {
	case "project":
		return c.Project.Get(strings.Join(nameParts[1:], "."))
	}

	return "", ErrKeyNotFound(name)
}
