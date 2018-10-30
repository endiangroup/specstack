package config

import "strings"

const (
	keyUser      prefix = "user"
	keyUserName         = "name"
	keyUserEmail        = "email"

	keyProject            prefix = "project"
	keyProjectRemote             = "remote"
	keyProjectName               = "name"
	keyProjectFeaturesDir        = "featuresdir"
	keyProjectPushingMode        = "pushingmode"
	keyProjectPullingMode        = "pullingmode"
)

func fetchPrefix(key string) prefix {
	keyParts := strings.Split(key, ".")

	return prefix(keyParts[0])
}

type prefix string

func (p prefix) Append(parts ...string) string {
	return strings.Join(append([]string{string(p)}, parts...), ".")
}

func fetchPostfix(key string) string {
	keyParts := strings.SplitN(key, ".", 2)

	return keyParts[len(keyParts)-1]
}
