package config

import (
	"github.com/endiangroup/specstack/errors"
	"github.com/endiangroup/specstack/validations"
)

type Validation func(*Config) error

func UserNameCannotBeBlank(c *Config) error {
	fieldName := KeyUser.Append(KeyUserName)
	if c.User == nil {
		return &errors.ValidationField{fieldName, "cannot be blank"}
	}

	return validations.CannotBeBlank(fieldName, c.User.Name)
}

func UserEmailCannotBeBlank(c *Config) error {
	fieldEmail := KeyUser.Append(KeyUserEmail)
	if c.User == nil {
		return &errors.ValidationField{fieldEmail, "cannot be blank"}
	}

	return validations.CannotBeBlank(fieldEmail, c.User.Email)
}
