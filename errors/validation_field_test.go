package errors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ValidationField(t *testing.T) {
	assert.EqualError(t, &ValidationField{}, "Field is invalid")

	assert.EqualError(t, &ValidationField{
		Field: "email",
	}, "Field 'email' is invalid")

	assert.EqualError(t, &ValidationField{
		Message: "cannot be blank",
	}, "Field cannot be blank")

	assert.EqualError(t, &ValidationField{
		Field:   "email",
		Message: "cannot be blank",
	}, "Field 'email' cannot be blank")
}
