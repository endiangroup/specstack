package config

import (
	"fmt"
	"strings"
)

type ErrKeyNotFound string

func (err ErrKeyNotFound) Error() string {
	return fmt.Sprintf("no config key '%s' found", string(err))
}

func Get(c *Config, name string) (string, error) {
	nameParts := strings.Split(name, ".")

	switch nameParts[0] {
	case "project":
		return projectGet(c.Project, strings.Join(nameParts[1:], "."))
	}

	return "", ErrKeyNotFound(name)
}

func projectGet(p *Project, name string) (string, error) {
	switch name {
	case "remote":
		return p.Remote, nil
	case "name":
		return p.Name, nil
	case "featuresdir":
		return p.FeaturesDir, nil
	case "pushingmode":
		return p.PushingMode, nil
	case "pullingmode":
		return p.PullingMode, nil
	}

	return "", ErrKeyNotFound(name)
}
