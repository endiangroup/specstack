package specification

type Snapshot struct {
	Scenarios []ScenarioSnapshot
}

type ScenarioSnapshot struct {
	StoryID    string
	LineNumber int
	ScenarioID string
}
