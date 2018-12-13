package specification

type Snapshotter struct {
	ReadSourcer  ReadSourcer
	ObjectHasher ObjectHasher
}

func (s *Snapshotter) Snapshot(spec *Specification) (Snapshot, error) {
	q := NewQuery(spec)
	q.MapReduce(MapScenarios()) //TODO: order by DID?

	snapshot := Snapshot{}
	for _, scenario := range q.Scenarios() {
		storyID, err := s.DeterministicID(scenario.Story)
		if err != nil {
			return snapshot, err
		}
		scenarioID, err := s.DeterministicID(scenario.Story)
		if err != nil {
			return snapshot, err
		}
		snapshot.Scenarios = append(snapshot.Scenarios, ScenarioSnapshot{
			StoryID:    storyID,
			ScenarioID: scenarioID,
			LineNumber: scenario.Location.Line,
		})
	}
	return snapshot, nil
}

func (s *Snapshotter) DeterministicID(object Sourcer) (string, error) {
	reader, err := s.ReadSourcer.ReadSource(object)
	if err != nil {
		return "", err
	}
	return s.ObjectHasher.ObjectHash(reader)
}
