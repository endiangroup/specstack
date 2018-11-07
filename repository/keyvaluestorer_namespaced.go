package repository

import "strings"

func NewNamespacedKeyValueStorer(keyValueStorer KeyValueStorer, namespace string) KeyValueStorer {
	return &namespacedKeyValueStorer{
		keyValueStorer: keyValueStorer,
		namespace:      namespace,
	}
}

type namespacedKeyValueStorer struct {
	keyValueStorer KeyValueStorer
	namespace      string
}

func (kv *namespacedKeyValueStorer) All() (map[string]string, error) {
	allKeyValues, err := kv.keyValueStorer.All()
	if err != nil {
		return nil, err
	}

	trimedKeyValues := map[string]string{}
	for key, value := range allKeyValues {
		if strings.HasPrefix(key, kv.namespace) {
			trimedKeyValues[kv.trimNamespace(key)] = value
		}
	}

	if len(trimedKeyValues) == 0 {
		return nil, GitConfigMissingKeyErr{}
	}

	return trimedKeyValues, nil
}

func (kv *namespacedKeyValueStorer) Get(key string) (string, error) {
	return kv.keyValueStorer.Get(kv.prefixNamespace(key))
}

func (kv *namespacedKeyValueStorer) Set(key, value string) error {
	return kv.keyValueStorer.Set(kv.prefixNamespace(key), value)
}

func (kv *namespacedKeyValueStorer) Unset(key string) error {
	return kv.keyValueStorer.Unset(kv.prefixNamespace(key))
}

func (kv *namespacedKeyValueStorer) prefixNamespace(key string) string {
	return kv.namespace + "." + key
}

func (kv *namespacedKeyValueStorer) trimNamespace(key string) string {
	return strings.TrimPrefix(key, kv.namespace+".")
}
