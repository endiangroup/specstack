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
	mockRepo := &repository.MockInitialiser{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	app := NewApp("", mockRepo, mockDeveloper, mockConfigStore)

	mockRepo.On("IsInitialised").Return(false)

	assert.Equal(t, ErrUninitialisedRepo, app.Initialise())
}

func Test_Initialise_CreatesConfigOnFirstRun(t *testing.T) {
	mockRepo := &repository.MockInitialiser{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	app := NewApp("", mockRepo, mockDeveloper, mockConfigStore)

	mockRepo.On("IsInitialised").Return(true)
	mockConfigStore.On("LoadConfig").Return(nil, persistence.ErrNoConfigFound)
	mockConfigStore.On("CreateConfig", mock.AnythingOfType("*config.Config")).Return(nil, nil)

	assert.NoError(t, app.Initialise())

	mockConfigStore.AssertExpectations(t)
}

func Test_Initialise_SetsConfigProjectNameToBaseOfPath(t *testing.T) {
	mockRepo := &repository.MockInitialiser{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	app := NewApp("/testing/test-dir", mockRepo, mockDeveloper, mockConfigStore)
	expectedConfig := config.NewWithDefaults()
	expectedConfig.Project.Name = "test-dir"

	mockRepo.On("IsInitialised").Return(true)
	mockConfigStore.On("LoadConfig").Return(nil, persistence.ErrNoConfigFound)
	mockConfigStore.On("CreateConfig", mock.Anything).Return(nil, nil)

	assert.NoError(t, app.Initialise())

	assert.Equal(t, expectedConfig, mockConfigStore.Calls[1].Arguments.Get(0))
}

func Test_Initialise_LoadsExistingConfigIfNotFirstRun(t *testing.T) {
	mockRepo := &repository.MockInitialiser{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	app := NewApp("/testing/test-dir", mockRepo, mockDeveloper, mockConfigStore)
	expectedConfig := config.NewWithDefaults()

	mockRepo.On("IsInitialised").Return(true)
	mockDeveloper.On("ListConfiguration", mock.Anything).Return(nil, nil)
	mockConfigStore.On("LoadConfig").Return(expectedConfig, nil)

	assert.NoError(t, app.Initialise())

	mockConfigStore.AssertExpectations(t)
}
