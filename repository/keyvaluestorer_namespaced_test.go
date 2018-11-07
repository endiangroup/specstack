package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get_PrefixesKeyWithNamespace(t *testing.T) {
	mockKVStorer := &MockConfigStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)

	mockKVStorer.On("GetConfig", namespace+".name").Return("", nil)

	_, err := kvstorer.GetConfig("name")
	assert.NoError(t, err)

	mockKVStorer.AssertExpectations(t)
}

func Test_Set_PrefixesKeyWithNamespace(t *testing.T) {
	mockKVStorer := &MockConfigStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)

	mockKVStorer.On("SetConfig", namespace+".name", "blah").Return(nil)

	assert.NoError(t, kvstorer.SetConfig("name", "blah"))

	mockKVStorer.AssertExpectations(t)
}

func Test_Unset_PrefixesKeyWithNamespace(t *testing.T) {
	mockKVStorer := &MockConfigStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)

	mockKVStorer.On("UnsetConfig", namespace+".name").Return(nil)

	assert.NoError(t, kvstorer.UnsetConfig("name"))

	mockKVStorer.AssertExpectations(t)
}

func Test_All_ReturnsOnlyConfigFromNamespaceAndTrims(t *testing.T) {
	mockKVStorer := &MockConfigStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)
	allKV := map[string]string{
		"testing.user.name":  "a b",
		"testing.user.email": "a@b.com",
		"unprefixed":         "123",
	}

	expectedKV := map[string]string{
		"user.name":  "a b",
		"user.email": "a@b.com",
	}

	mockKVStorer.On("AllConfig").Return(allKV, nil)

	returnedKV, err := kvstorer.AllConfig()
	assert.NoError(t, err)

	assert.Equal(t, expectedKV, returnedKV)
}

func Test_All_ReturnsKeyMissingErrorIfNoKeysInNamespace(t *testing.T) {
	mockKVStorer := &MockConfigStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)
	allKV := map[string]string{
		"unprefixed": "123",
	}

	mockKVStorer.On("AllConfig").Return(allKV, nil)

	_, err := kvstorer.AllConfig()

	assert.Equal(t, GitConfigMissingKeyErr{}, err)
}
