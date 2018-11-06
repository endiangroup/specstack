package cmd

import (
	"strings"

	"github.com/endiangroup/specstack/errors"
)

type Validation func(string) error

func IsKeyEqualsValueFormat(arg string) error {
	parts := strings.Split(arg, "=")
	if len(parts) != 2 {
		return &errors.ValidationError{E: errors.New("invalid argument format, expected: key=value")}
	}

	return nil
}
