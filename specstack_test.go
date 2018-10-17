package specstack

import (
	"testing"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/personas"
	"github.com/endiangroup/specstack/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Initialise_ReturnsErrorIfRepoisotiryIsntInitialised(t *testing.T) {
	mockRepo := &repository.MockReadWriter{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	app := NewApp("", mockRepo, mockDeveloper, mockConfigStore)

	mockRepo.On("IsInitialised").Return(false)

	assert.Equal(t, ErrUninitialisedRepo, app.Initialise())
}

func Test_Initialise_CreatesConfigOnFirstRun(t *testing.T) {
	mockRepo := &repository.MockReadWriter{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	app := NewApp("", mockRepo, mockDeveloper, mockConfigStore)

	mockRepo.On("IsInitialised").Return(true)
	mockDeveloper.On("ListConfiguration").Return(nil, persistence.ErrNoConfigFound)
	mockConfigStore.On("CreateConfig", mock.Anything).Return(nil, nil)

	assert.NoError(t, app.Initialise())

	mockConfigStore.AssertExpectations(t)
}

func Test_Initialise_SetsConfigProjectNameToBaseOfPath(t *testing.T) {
	mockRepo := &repository.MockReadWriter{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	app := NewApp("/testing/test-dir", mockRepo, mockDeveloper, mockConfigStore)
	expectedConfig := config.NewWithDefaults()
	expectedConfig.Project.Name = "test-dir"

	mockRepo.On("IsInitialised").Return(true)
	mockDeveloper.On("ListConfiguration").Return(nil, persistence.ErrNoConfigFound)
	mockConfigStore.On("CreateConfig", mock.Anything).Return(nil, nil)

	assert.NoError(t, app.Initialise())

	assert.Equal(t, expectedConfig, mockConfigStore.Calls[0].Arguments.Get(0))
}
