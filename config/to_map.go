package config

func ToMap(c *Config) map[string]string {
	return projectToMap(c.Project, "project.", map[string]string{})
}

func projectToMap(p *Project, prefix string, configMap map[string]string) map[string]string {
	configMap[prefix+"remote"] = p.Remote
	configMap[prefix+"name"] = p.Name
	configMap[prefix+"featuresdir"] = p.FeaturesDir
	configMap[prefix+"pushingmode"] = p.PushingMode
	configMap[prefix+"pullingmode"] = p.PullingMode

	return configMap
}
