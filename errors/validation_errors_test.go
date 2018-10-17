package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidationErrors_RendersSeveralErrors(t *testing.T) {
	assert.EqualError(
		t,
		ValidationErrors{&ValidationField{"a", "is invalid"}, &ValidationField{"b", "is invalid"}},
		"Field 'a' is invalid, Field 'b' is invalid",
	)
}

func Test_ValidationErrors_AppendsErrors(t *testing.T) {
	subject0 := ValidationErrors{}
	assert.EqualError(t, subject0, "")

	subject1 := subject0.Append(&ValidationField{"a", "is invalid"})
	assert.EqualError(t, subject0, "")
	assert.EqualError(t, subject1, "Field 'a' is invalid")

	subject2 := subject1.Append(&ValidationField{"b", "is invalid"})
	assert.EqualError(t, subject0, "")
	assert.EqualError(t, subject1, "Field 'a' is invalid")
	assert.EqualError(t, subject2, "Field 'a' is invalid, Field 'b' is invalid")
}

func Test_ValidationErrors_Any(t *testing.T) {
	subject0 := ValidationErrors{}
	assert.False(t, subject0.Any())

	subject1 := subject0.Append(&ValidationField{"a", "is invalid"})
	assert.True(t, subject1.Any())
}
