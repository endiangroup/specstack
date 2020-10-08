package config

const (
	ModeAuto     = "auto"
	ModeSemiAuto = "semi-auto"
)

func newProject() *Project {
	return &Project{}
}

func newProjectWithDefaults() *Project {
	return &Project{
		Remote:      "origin",
		FeaturesDir: "./features",
		PushingMode: ModeAuto,
		PullingMode: ModeSemiAuto,
	}
}

type Project struct {
	Remote      string
	Name        string
	FeaturesDir string
	PushingMode string
	PullingMode string
}
