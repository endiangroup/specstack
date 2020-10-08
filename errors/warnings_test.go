package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Warnings_RendersSeveralErrors(t *testing.T) {
	assert.EqualError(
		t,
		NewWarnings(errors.New("a"), errors.New("b")),
		"a, b",
	)
}

func Test_Warnings_AppendsErrors(t *testing.T) {
	subject0 := Warnings{}
	assert.EqualError(t, subject0, "")

	subject1 := subject0.Append(errors.New("a"))
	assert.EqualError(t, subject0, "")
	assert.EqualError(t, subject1, "a")

	subject2 := subject1.Append(errors.New("b"))
	assert.EqualError(t, subject0, "")
	assert.EqualError(t, subject1, "a")
	assert.EqualError(t, subject2, "a, b")
}

func Test_Warnings_Any(t *testing.T) {
	subject0 := Warnings{}
	assert.False(t, subject0.Any())

	subject1 := subject0.Append(errors.New("a"))
	assert.True(t, subject1.Any())
}

func Test_WarningOrNil_ProducesCorrectResults(t *testing.T) {
	assert.Nil(t, WarningOrNil(nil))
	assert.Nil(t, WarningOrNil(NewWarnings()))
	assert.Equal(t, WarningOrNil(errors.New("a")), NewWarning(errors.New("a")))
	assert.EqualError(t, WarningOrNil(NewWarnings(errors.New("a"), errors.New("b"))), "a, b")
}

func Test_IsWarning_ProducesCorrectResults(t *testing.T) {
	assert.False(t, IsWarning(nil))
	assert.False(t, IsWarning(errors.New("a")))
	assert.False(t, IsWarning(WarningOrNil(NewWarnings())))
	assert.True(t, IsWarning(NewWarning(errors.New("a"))))
	assert.True(t, IsWarning(NewWarnings(errors.New("a"), errors.New("b"))))
}
