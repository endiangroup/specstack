package config

import "github.com/endiangroup/specstack/errors"

func IsValid(c *Config, validations ...Validation) (bool, error) {
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
