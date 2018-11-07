package repository

import "strings"

func NewNamespacedKeyValueStorer(keyValueStorer KeyValueStorer, namespace string) KeyValueStorer {
	return &NamespacedKeyValueStorer{
		keyValueStorer: keyValueStorer,
		namespace:      namespace,
	}
}

type NamespacedKeyValueStorer struct {
	keyValueStorer KeyValueStorer
	namespace      string
}

func (kv *NamespacedKeyValueStorer) All() (map[string]string, error) {
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

func (kv *NamespacedKeyValueStorer) Get(key string) (string, error) {
	return kv.keyValueStorer.Get(kv.prefixNamespace(key))
}

func (kv *NamespacedKeyValueStorer) Set(key, value string) error {
	return kv.keyValueStorer.Set(kv.prefixNamespace(key), value)
}

func (kv *NamespacedKeyValueStorer) Unset(key string) error {
	return kv.keyValueStorer.Unset(kv.prefixNamespace(key))
}

func (kv *NamespacedKeyValueStorer) prefixNamespace(key string) string {
	return kv.namespace + "." + key
}

func (kv *NamespacedKeyValueStorer) trimNamespace(key string) string {
	return strings.TrimPrefix(key, kv.namespace+".")
}
