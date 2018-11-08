package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get_PrefixesKeyWithNamespace(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	namespace := "testing"
	configstorer := NewNamespacedKeyValueStorer(mockConfigStorer, namespace)

	mockConfigStorer.On("GetConfig", namespace+".name").Return("", nil)

	_, err := configstorer.GetConfig("name")
	assert.NoError(t, err)

	mockConfigStorer.AssertExpectations(t)
}

func Test_Set_PrefixesKeyWithNamespace(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	namespace := "testing"
	configstorer := NewNamespacedKeyValueStorer(mockConfigStorer, namespace)

	mockConfigStorer.On("SetConfig", namespace+".name", "blah").Return(nil)

	assert.NoError(t, configstorer.SetConfig("name", "blah"))

	mockConfigStorer.AssertExpectations(t)
}

func Test_Unset_PrefixesKeyWithNamespace(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	namespace := "testing"
	configstorer := NewNamespacedKeyValueStorer(mockConfigStorer, namespace)

	mockConfigStorer.On("UnsetConfig", namespace+".name").Return(nil)

	assert.NoError(t, configstorer.UnsetConfig("name"))

	mockConfigStorer.AssertExpectations(t)
}

func Test_All_ReturnsOnlyConfigFromNamespaceAndTrims(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	namespace := "testing"
	configstorer := NewNamespacedKeyValueStorer(mockConfigStorer, namespace)
	allKV := map[string]string{
		"testing.user.name":  "a b",
		"testing.user.email": "a@b.com",
		"unprefixed":         "123",
	}

	expectedKV := map[string]string{
		"user.name":  "a b",
		"user.email": "a@b.com",
	}

	mockConfigStorer.On("AllConfig").Return(allKV, nil)

	returnedKV, err := configstorer.AllConfig()
	assert.NoError(t, err)

	assert.Equal(t, expectedKV, returnedKV)
}

func Test_All_ReturnsKeyMissingErrorIfNoKeysInNamespace(t *testing.T) {
	mockConfigStorer := &MockConfigStorer{}
	namespace := "testing"
	configstorer := NewNamespacedKeyValueStorer(mockConfigStorer, namespace)
	allKV := map[string]string{
		"unprefixed": "123",
	}

	mockConfigStorer.On("AllConfig").Return(allKV, nil)

	_, err := configstorer.AllConfig()

	assert.Equal(t, GitConfigMissingKeyErr{}, err)
}
