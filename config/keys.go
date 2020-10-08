package config

import "strings"

const (
	KeyUser      prefix = "user"
	KeyUserName         = "name"
	KeyUserEmail        = "email"

	KeyProject            prefix = "project"
	KeyProjectRemote             = "remote"
	KeyProjectName               = "name"
	KeyProjectFeaturesDir        = "featuresdir"
	KeyProjectPushingMode        = "pushingmode"
	KeyProjectPullingMode        = "pullingmode"
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
