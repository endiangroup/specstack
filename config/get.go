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
	case KeyProject:
		return projectGet(c.Project, key)
	case KeyUser:
		return userGet(c.User, key)
	}

	return "", ErrKeyNotFound(key)
}

func userGet(u *User, key string) (string, error) {
	switch fetchPostfix(key) {
	case KeyUserName:
		return u.Name, nil
	case KeyUserEmail:
		return u.Email, nil
	}

	return "", ErrKeyNotFound(key)
}

func projectGet(p *Project, key string) (string, error) {
	switch fetchPostfix(key) {
	case KeyProjectName:
		return p.Name, nil
	case KeyProjectRemote:
		return p.Remote, nil
	case KeyProjectFeaturesDir:
		return p.FeaturesDir, nil
	case KeyProjectPushingMode:
		return p.PushingMode, nil
	case KeyProjectPullingMode:
		return p.PullingMode, nil
	}

	return "", ErrKeyNotFound(key)
}
