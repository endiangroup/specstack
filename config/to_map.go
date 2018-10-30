package config

func ToMap(c *Config) map[string]string {
	target := projectToMap(c.Project, keyProject, map[string]string{})

	return userToMap(c.User, keyUser, target)
}

func userToMap(u *User, key prefix, configMap map[string]string) map[string]string {
	configMap[key.Append(keyUserName)] = u.Name
	configMap[key.Append(keyUserEmail)] = u.Email

	return configMap
}

func projectToMap(p *Project, key prefix, configMap map[string]string) map[string]string {
	configMap[key.Append(keyProjectName)] = p.Name
	configMap[key.Append(keyProjectRemote)] = p.Remote
	configMap[key.Append(keyProjectFeaturesDir)] = p.FeaturesDir
	configMap[key.Append(keyProjectPushingMode)] = p.PushingMode
	configMap[key.Append(keyProjectPullingMode)] = p.PullingMode

	return configMap
}
