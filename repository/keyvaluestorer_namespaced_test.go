package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Get_PrefixesKeyWithNamespace(t *testing.T) {
	mockKVStorer := &MockKeyValueStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)

	mockKVStorer.On("Get", namespace+".name").Return("", nil)

	_, err := kvstorer.Get("name")
	assert.NoError(t, err)

	mockKVStorer.AssertExpectations(t)
}

func Test_Set_PrefixesKeyWithNamespace(t *testing.T) {
	mockKVStorer := &MockKeyValueStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)

	mockKVStorer.On("Set", namespace+".name", "blah").Return(nil)

	assert.NoError(t, kvstorer.Set("name", "blah"))

	mockKVStorer.AssertExpectations(t)
}

func Test_Unset_PrefixesKeyWithNamespace(t *testing.T) {
	mockKVStorer := &MockKeyValueStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)

	mockKVStorer.On("Unset", namespace+".name").Return(nil)

	assert.NoError(t, kvstorer.Unset("name"))

	mockKVStorer.AssertExpectations(t)
}

func Test_All_ReturnsOnlyConfigFromNamespaceAndTrims(t *testing.T) {
	mockKVStorer := &MockKeyValueStorer{}
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

	mockKVStorer.On("All").Return(allKV, nil)

	returnedKV, err := kvstorer.All()
	assert.NoError(t, err)

	assert.Equal(t, expectedKV, returnedKV)
}

func Test_All_ReturnsKeyMissingErrorIfNoKeysInNamespace(t *testing.T) {
	mockKVStorer := &MockKeyValueStorer{}
	namespace := "testing"
	kvstorer := NewNamespacedKeyValueStorer(mockKVStorer, namespace)
	allKV := map[string]string{
		"unprefixed": "123",
	}

	mockKVStorer.On("All").Return(allKV, nil)

	_, err := kvstorer.All()

	assert.Equal(t, GitConfigMissingKeyErr{}, err)
}
