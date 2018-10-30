package persistence

import (
	"testing"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/errors"
	"github.com/endiangroup/specstack/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_StoreConfig_SetsAllConfigKeyValuesOnRepository(t *testing.T) {
	mockKVStore := &repository.MockKeyValueStorer{}
	repoStore := NewRepositoryStore(mockKVStore)
	c := config.NewWithDefaults()

	mockKVStore.On("Set", mock.Anything, mock.Anything).Return(nil)

	_, err := repoStore.StoreConfig(c)
	assert.NoError(t, err)

	configMap := config.ToMap(c)
	for _, call := range mockKVStore.Calls {
		key := call.Arguments.String(0)
		assert.Contains(t, configMap, key)

		assert.Equal(t, configMap[key], call.Arguments.String(1))
	}
}

func Test_StoreConfig_ReturnsAnyConfigSetErrors(t *testing.T) {
	mockKVStore := &repository.MockKeyValueStorer{}
	repoStore := NewRepositoryStore(mockKVStore)
	config := config.NewWithDefaults()

	mockKVStore.On("Set", mock.Anything, mock.Anything).Return(errors.New("!!!"))

	_, err := repoStore.StoreConfig(config)

	assert.True(t, len(err.(errors.Errors)) > 0)
}

func Test_LoadConfig_SetsKeyValuesOnConfig(t *testing.T) {
	mockKVStore := &repository.MockKeyValueStorer{}
	repoStore := NewRepositoryStore(mockKVStore)
	c := config.New()
	c.Project.Remote = "upstream"
	c.Project.Name = "test"
	expectedConfigMap := config.ToMap(c)

	mockKVStore.On("All").Return(expectedConfigMap, nil)

	c, err := repoStore.LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, config.ToMap(c), expectedConfigMap)
}
