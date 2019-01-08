package personas

import (
	"testing"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/repository"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

func Test_DeveloperAssertConfig_CreatesConfigOnFirstRun(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	dev := &Developer{
		repo:  mockRepo,
		store: repoStore,
	}

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("user@email", nil)
	mockRepo.On("PrepareMetadataSync").Return(nil)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, nil)

	assert.NoError(t, dev.AssertConfig())

	mockConfigStore.AssertExpectations(t)
}

func Test_DeveloperAssertConfig_ReturnsErrorWhenMissingUsername(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	dev := &Developer{
		repo:  mockRepo,
		store: repoStore,
	}

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("", persistence.ErrNoConfigFound)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, persistence.ErrNoConfigFound)

	err := dev.AssertConfig()

	assert.IsType(t, MissingRequiredConfigValueErr(""), err)
}

func Test_DeveloperAssertConfig_ReturnsErrorWhenMissingEmail(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})
	dev := &Developer{
		repo:  mockRepo,
		store: repoStore,
	}

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("", persistence.ErrNoConfigFound)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, persistence.ErrNoConfigFound)

	err := dev.AssertConfig()

	assert.IsType(t, MissingRequiredConfigValueErr(""), err)
}

func Test_DeveloperAssertConfig_SetsConfigDefaults(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})

	dev := &Developer{
		path:  "test-dir",
		repo:  mockRepo,
		store: repoStore,
	}

	expectedConfig := config.NewWithDefaults()
	expectedConfig.Project.Name = "test-dir"
	expectedConfig.User.Name = "username"
	expectedConfig.User.Email = "user@email"

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("GetConfig", "user.name").Return("username", nil)
	mockRepo.On("GetConfig", "user.email").Return("user@email", nil)
	mockRepo.On("PrepareMetadataSync").Return(nil)
	mockConfigStore.On("SetConfig", mock.Anything, mock.Anything).Return(nil)
	mockConfigStore.On("AllConfig").Return(map[string]string{}, persistence.ErrNoConfigFound)
	mockConfigStore.On("StoreConfig", mock.Anything).Return(nil, nil)

	assert.NoError(t, dev.AssertConfig())

	assert.Equal(t, expectedConfig, dev.config)
}

func Test_DeveloperAssertConfig_LoadsExistingConfigIfNotFirstRun(t *testing.T) {
	mockRepo := &repository.MockRepository{}
	mockConfigStore := &persistence.MockConfigStorer{}
	repoStore := persistence.NewStore(mockConfigStore, &persistence.MockMetadataStorer{})

	dev := &Developer{
		path:  "test-dir",
		repo:  mockRepo,
		store: repoStore,
	}
	expectedConfig := config.NewWithDefaults()

	mockRepo.On("IsInitialised").Return(true)
	mockRepo.On("PrepareMetadataSync").Return(nil)
	mockConfigStore.On("AllConfig").Return(config.ToMap(expectedConfig), nil)

	assert.NoError(t, dev.AssertConfig())

	mockConfigStore.AssertExpectations(t)
	assert.Equal(t, expectedConfig, dev.config)
}
