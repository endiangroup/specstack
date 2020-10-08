package validations

import (
	"strings"

	"github.com/endiangroup/specstack/errors"
)

func CannotBeBlank(field, value string) error {
	if strings.TrimSpace(value) == "" {
		return &errors.ValidationField{Field: field, Message: "cannot be blank"}
	}

	return nil
}
