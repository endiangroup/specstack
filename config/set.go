package config

import "strings"

func Set(c *Config, name, value string) {
	nameParts := strings.Split(name, ".")

	switch nameParts[0] {
	case "project":
		projectSet(c.Project, strings.Join(nameParts[1:], "."), value)
	}
}

func projectSet(p *Project, name, value string) {
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
	}
}
