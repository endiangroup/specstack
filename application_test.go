package specstack

import (
	"errors"
	"testing"

	"github.com/endiangroup/specstack/repository"
	"github.com/stretchr/testify/assert"
)

func Test_Initialise_ReturnsErrorIfRepositoryIsntInitialised(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	app := &Application{
		Repository: mockRepo,
	}

	mockRepo.On("IsInitialised").Return(false)

	assert.Equal(t, ErrUninitialisedRepo, app.Initialise())
}

func Test_Initialise_ReturnsErrorIfConfigAssertionFails(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockConfigAsserter := &MockConfigAsserter{}
	app := &Application{
		Repository:     mockRepo,
		ConfigAsserter: mockConfigAsserter,
	}

	mockRepo.On("IsInitialised").Return(true)
	mockConfigAsserter.On("AssertConfig").Return(errors.New("!!!"))

	assert.Equal(t, errors.New("!!!"), app.Initialise())
}
