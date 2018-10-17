package errors

import (
	"errors"
	"strings"
)

var New = errors.New

type Errors []error

func (errs Errors) Error() string {
	components := []string{}
	for _, err := range errs {
		components = append(components, err.Error())
	}
	return strings.Join(components, ", ")
}

func (errs Errors) Append(err error) Errors {
	return append(errs, err)
}

func (errs Errors) Any() bool {
	return len(errs) > 0
}
