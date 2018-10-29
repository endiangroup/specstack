package config

func ToMap(c *Config) map[string]string {
	return projectToMap(c.Project, keyProject, map[string]string{})
}

func projectToMap(p *Project, key prefix, configMap map[string]string) map[string]string {
	configMap[key.Append(keyProjectName)] = p.Name
	configMap[key.Append(keyProjectRemote)] = p.Remote
	configMap[key.Append(keyProjectFeaturesDir)] = p.FeaturesDir
	configMap[key.Append(keyProjectPushingMode)] = p.PushingMode
	configMap[key.Append(keyProjectPullingMode)] = p.PullingMode

	return configMap
}
