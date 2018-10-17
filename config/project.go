package config

func newProject() *Project {
	return &Project{}
}

func NewProjectWithDefaults() *Project {
	return &Project{
		Remote:      "origin",
		FeaturesDir: "./features",
		PushingMode: "auto",
		PullingMode: "semi-auto",
	}
}

type Project struct {
	Remote      string
	Name        string
	FeaturesDir string
	PushingMode string
	PullingMode string
}

func (p *Project) ToMap(prefix string, configMap map[string]string) map[string]string {
	configMap[prefix+"remote"] = p.Remote
	configMap[prefix+"name"] = p.Name
	configMap[prefix+"featuresdir"] = p.FeaturesDir
	configMap[prefix+"pushingmode"] = p.PushingMode
	configMap[prefix+"pullingmode"] = p.PullingMode

	return configMap
}

func (p *Project) Set(name, value string) {
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
