package config

import "strings"

func Set(c *Config, name, value string) error {
	nameParts := strings.Split(name, ".")

	switch nameParts[0] {
	case "project":
		return projectSet(c.Project, strings.Join(nameParts[1:], "."), value)
	}

	return ErrKeyNotFound(name)
}

func projectSet(p *Project, name, value string) error {
	switch name {
	case "remote":
		p.Remote = value
	case "name":
		p.Name = value
	case "featuresdir":
		p.FeaturesDir = value
	case "pushingmode":
		p.PushingMode = value
	case "pullingmode":
		p.PullingMode = value
	default:
		return ErrKeyNotFound(name)
	}

	return nil
}
