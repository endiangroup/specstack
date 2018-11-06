package config

func Set(c *Config, key, value string) error {
	switch fetchPrefix(key) {
	case keyProject:
		return projectSet(c.Project, key, value)
	case keyUser:
		return userSet(c.User, key, value)
	}

	return ErrKeyNotFound(key)
}

func userSet(u *User, key, value string) error {
	switch fetchPostfix(key) {
	case keyUserName:
		u.Name = value
	case keyUserEmail:
		u.Email = value
	default:
		return ErrKeyNotFound(key)
	}

	return nil
}

func projectSet(p *Project, key, value string) error {
	switch fetchPostfix(key) {
	case keyProjectName:
		p.Name = value
	case keyProjectRemote:
		p.Remote = value
	case keyProjectFeaturesDir:
		p.FeaturesDir = value
	case keyProjectPushingMode:
		p.PushingMode = value
	case keyProjectPullingMode:
		p.PullingMode = value
	default:
		return ErrKeyNotFound(key)
	}

	return nil
}
