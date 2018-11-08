package config

func Set(c *Config, key, value string) error {
	switch fetchPrefix(key) {
	case KeyProject:
		return projectSet(c.Project, key, value)
	case KeyUser:
		return userSet(c.User, key, value)
	}

	return ErrKeyNotFound(key)
}

func userSet(u *User, key, value string) error {
	switch fetchPostfix(key) {
	case KeyUserName:
		u.Name = value
	case KeyUserEmail:
		u.Email = value
	default:
		return ErrKeyNotFound(key)
	}

	return nil
}

func projectSet(p *Project, key, value string) error {
	switch fetchPostfix(key) {
	case KeyProjectName:
		p.Name = value
	case KeyProjectRemote:
		p.Remote = value
	case KeyProjectFeaturesDir:
		p.FeaturesDir = value
	case KeyProjectPushingMode:
		p.PushingMode = value
	case KeyProjectPullingMode:
		p.PullingMode = value
	default:
		return ErrKeyNotFound(key)
	}

	return nil
}
