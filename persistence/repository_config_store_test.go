package persistence

import (
	"testing"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/errors"
	"github.com/endiangroup/specstack/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_CreateConfig_SetsAllConfigKeyValuesOnRepository(t *testing.T) {
	mockRepo := &repository.MockReadWriter{}
	repoStore := NewRepositoryStore(mockRepo)
	c := config.NewWithDefaults()

	mockRepo.On("ConfigSet", mock.Anything, mock.Anything).Return(nil)

	_, err := repoStore.CreateConfig(c)
	assert.NoError(t, err)

	configMap := config.ToMap(c)
	for _, call := range mockRepo.Calls {
		key := call.Arguments.String(0)
		assert.Contains(t, configMap, key)

		assert.Equal(t, configMap[key], call.Arguments.String(1))
	}
}

func Test_CreateConfig_ReturnsAnyConfigSetErrors(t *testing.T) {
	mockRepo := &repository.MockReadWriter{}
	repoStore := NewRepositoryStore(mockRepo)
	config := config.NewWithDefaults()

	mockRepo.On("ConfigSet", mock.Anything, mock.Anything).Return(errors.New("!!!"))

	_, err := repoStore.CreateConfig(config)

	assert.True(t, len(err.(errors.Errors)) > 0)
}

func Test_LoadConfig_SetsKeyValuesOnConfig(t *testing.T) {
	mockRepo := &repository.MockReadWriter{}
	repoStore := NewRepositoryStore(mockRepo)
	expectedConfigMap := map[string]string{
		"project.featuresdir": "",
		"project.pushingmode": "",
		"project.pullingmode": "",
		"project.remote":      "upstream",
		"project.name":        "test",
	}

	mockRepo.On("ConfigGetAll").Return(expectedConfigMap, nil)

	c, err := repoStore.LoadConfig()
	assert.NoError(t, err)

	assert.Equal(t, config.ToMap(c), expectedConfigMap)
}
