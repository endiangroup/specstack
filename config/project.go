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

func (p *Project) Get(name string) (string, error) {
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
