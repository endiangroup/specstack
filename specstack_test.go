package specstack

import (
	"testing"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/metadata"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/personas"
	"github.com/endiangroup/specstack/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_Initialise_ReturnsErrorIfRepoisotiryIsntInitialised(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	mockMetadataStore := &metadata.MockReadStorer{}
	app := New("", mockRepo, mockDeveloper, mockConfigStore, mockMetadataStore)

	mockRepo.On("IsInitialised").Return(false)

	assert.Equal(t, ErrUninitialisedRepo, app.Initialise())
}

func Test_Initialise_CreatesConfigOnFirstRun(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	mockMetadataStore := &metadata.MockReadStorer{}
	app := New("", mockRepo, mockDeveloper, mockConfigStore, mockMetadataStore)

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("user@email", nil)
	mockConfigStore.On("LoadConfig").Return(nil, persistence.ErrNoConfigFound)
	mockConfigStore.On("StoreConfig", mock.AnythingOfType("*config.Config")).Return(nil, nil)

	assert.NoError(t, app.Initialise())

	mockConfigStore.AssertExpectations(t)
}

func Test_Initialise_ReturnsErrorWhenMissingUsername(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	mockMetadataStore := &metadata.MockReadStorer{}
	app := New("/testing/test-dir", mockRepo, mockDeveloper, mockConfigStore, mockMetadataStore)

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("", persistence.ErrNoConfigFound)
	mockConfigStore.On("LoadConfig").Return(nil, persistence.ErrNoConfigFound)

	err := app.Initialise()

	assert.IsType(t, MissingRequiredConfigValueErr(""), err)
}

func Test_Initialise_ReturnsErrorWhenMissingEmail(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	mockMetadataStore := &metadata.MockReadStorer{}
	app := New("/testing/test-dir", mockRepo, mockDeveloper, mockConfigStore, mockMetadataStore)

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("", persistence.ErrNoConfigFound)
	mockConfigStore.On("LoadConfig").Return(nil, persistence.ErrNoConfigFound)

	err := app.Initialise()

	assert.IsType(t, MissingRequiredConfigValueErr(""), err)
}

func Test_Initialise_SetsConfigDefaults(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	mockMetadataStore := &metadata.MockReadStorer{}
	app := New("/testing/test-dir", mockRepo, mockDeveloper, mockConfigStore, mockMetadataStore)
	expectedConfig := config.NewWithDefaults()
	expectedConfig.Project.Name = "test-dir"
	expectedConfig.User.Name = "username"
	expectedConfig.User.Email = "user@email"

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("user@email", nil)
	mockConfigStore.On("LoadConfig").Return(nil, persistence.ErrNoConfigFound)
	mockConfigStore.On("StoreConfig", mock.Anything).Return(nil, nil)

	assert.NoError(t, app.Initialise())

	assert.Equal(t, expectedConfig, mockConfigStore.Calls[1].Arguments.Get(0))
}

func Test_Initialise_LoadsExistingConfigIfNotFirstRun(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockDeveloper := &personas.MockDeveloper{}
	mockConfigStore := &config.MockStorer{}
	mockMetadataStore := &metadata.MockReadStorer{}
	app := New("/testing/test-dir", mockRepo, mockDeveloper, mockConfigStore, mockMetadataStore)
	expectedConfig := config.NewWithDefaults()

	mockRepo.On("IsInitialised").Return(true)
	mockDeveloper.On("ListConfiguration", mock.Anything).Return(nil, nil)
	mockConfigStore.On("LoadConfig").Return(expectedConfig, nil)

	assert.NoError(t, app.Initialise())

	mockConfigStore.AssertExpectations(t)
}
