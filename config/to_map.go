package config

func ToMap(c *Config) map[string]string {
	target := projectToMap(c.Project, KeyProject, map[string]string{})

	return userToMap(c.User, KeyUser, target)
}

func userToMap(u *User, key prefix, configMap map[string]string) map[string]string {
	configMap[key.Append(KeyUserName)] = u.Name
	configMap[key.Append(KeyUserEmail)] = u.Email

	return configMap
}

func projectToMap(p *Project, key prefix, configMap map[string]string) map[string]string {
	configMap[key.Append(KeyProjectName)] = p.Name
	configMap[key.Append(KeyProjectRemote)] = p.Remote
	configMap[key.Append(KeyProjectFeaturesDir)] = p.FeaturesDir
	configMap[key.Append(KeyProjectPushingMode)] = p.PushingMode
	configMap[key.Append(KeyProjectPullingMode)] = p.PullingMode

	return configMap
}
