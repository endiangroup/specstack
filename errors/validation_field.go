package errors

import (
	"fmt"
)

type ValidationField struct {
	Field   string
	Message string
}

func (err *ValidationField) Error() string {
	message := err.Message
	if message == "" {
		message = "is invalid"
	}

	if err.Field != "" {
		return fmt.Sprintf("Field '%s' %s", err.Field, message)
	}
	return fmt.Sprintf("Field %s", message)
}
