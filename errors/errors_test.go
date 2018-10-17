package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Errors_RendersSeveralErrors(t *testing.T) {
	assert.EqualError(
		t,
		Errors{errors.New("a"), errors.New("b")},
		"a, b",
	)
}

func Test_Errors_AppendsErrors(t *testing.T) {
	subject0 := Errors{}
	assert.EqualError(t, subject0, "")

	subject1 := subject0.Append(errors.New("a"))
	assert.EqualError(t, subject0, "")
	assert.EqualError(t, subject1, "a")

	subject2 := subject1.Append(errors.New("b"))
	assert.EqualError(t, subject0, "")
	assert.EqualError(t, subject1, "a")
	assert.EqualError(t, subject2, "a, b")
}

func Test_Errors_Any(t *testing.T) {
	subject0 := Errors{}
	assert.False(t, subject0.Any())

	subject1 := subject0.Append(errors.New("a"))
	assert.True(t, subject1.Any())
}
