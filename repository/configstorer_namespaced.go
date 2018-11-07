package repository

import "strings"

func NewNamespacedKeyValueStorer(keyValueStorer ConfigStorer, namespace string) ConfigStorer {
	return &namespacedKeyValueStorer{
		keyValueStorer: keyValueStorer,
		namespace:      namespace,
	}
}

type namespacedKeyValueStorer struct {
	keyValueStorer ConfigStorer
	namespace      string
}

func (kv *namespacedKeyValueStorer) AllConfig() (map[string]string, error) {
	allKeyValues, err := kv.keyValueStorer.AllConfig()
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

func (kv *namespacedKeyValueStorer) GetConfig(key string) (string, error) {
	return kv.keyValueStorer.GetConfig(kv.prefixNamespace(key))
}

func (kv *namespacedKeyValueStorer) SetConfig(key, value string) error {
	return kv.keyValueStorer.SetConfig(kv.prefixNamespace(key), value)
}

func (kv *namespacedKeyValueStorer) UnsetConfig(key string) error {
	return kv.keyValueStorer.UnsetConfig(kv.prefixNamespace(key))
}

func (kv *namespacedKeyValueStorer) prefixNamespace(key string) string {
	return kv.namespace + "." + key
}

func (kv *namespacedKeyValueStorer) trimNamespace(key string) string {
	return strings.TrimPrefix(key, kv.namespace+".")
}
