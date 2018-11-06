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
	case keyUser:
		return userGet(c.User, key)
	}

	return "", ErrKeyNotFound(key)
}

func userGet(u *User, key string) (string, error) {
	switch fetchPostfix(key) {
	case keyUserName:
		return u.Name, nil
	case keyUserEmail:
		return u.Email, nil
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
