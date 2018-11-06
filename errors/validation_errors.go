package errors

import (
	"strings"
)

type ValidationErrors Errors

func (errs ValidationErrors) Error() string {
	components := []string{}
	for _, err := range errs {
		components = append(components, err.Error())
	}
	return strings.Join(components, ", ")
}

func (errs ValidationErrors) Append(err error) ValidationErrors {
	return append(errs, err)
}

func (errs ValidationErrors) Any() bool {
	return Errors(errs).Any()
}
