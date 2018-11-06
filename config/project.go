package config

func newProject() *Project {
	return &Project{}
}

func newProjectWithDefaults() *Project {
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
