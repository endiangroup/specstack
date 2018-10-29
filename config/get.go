package config

import (
	"fmt"
)

type ErrKeyNotFound string

func (err ErrKeyNotFound) Error() string {
	return fmt.Sprintf("no config key '%s' found", string(err))
}

func Get(c *Config, key string) (string, error) {
	switch fetchPrefix(key) {
	case keyProject:
		return projectGet(c.Project, key)
	}

	return "", ErrKeyNotFound(key)
}

func projectGet(p *Project, key string) (string, error) {
	switch fetchPostfix(key) {
	case keyProjectName:
		return p.Name, nil
	case keyProjectRemote:
		return p.Remote, nil
	case keyProjectFeaturesDir:
		return p.FeaturesDir, nil
	case keyProjectPushingMode:
		return p.PushingMode, nil
	case keyProjectPullingMode:
		return p.PullingMode, nil
	}

	return "", ErrKeyNotFound(key)
}
