package specification

import "reflect"

type Snapshot struct {
	Scenarios []ScenarioSnapshot
}

type ScenarioSnapshot struct {
	StorySource Source
	StoryID     string
	LineNumber  int
	ScenarioID  string
}

func (a Snapshot) Equal(b Snapshot) bool {
	return reflect.DeepEqual(a, b)
}

/*
Diff returns the differences between Snapshots A and B. It returns
two new snapshots: one that includes all the elements that are in A
but not in B, and another that includes all elements that are in B
but not in A.
*/
func (a Snapshot) Diff(b Snapshot) (removed, added Snapshot) {
	removed.Scenarios = a.diffScenarios(a.Scenarios, b.Scenarios)
	added.Scenarios = a.diffScenarios(b.Scenarios, a.Scenarios)
	return
}

func (s Snapshot) diffScenarios(a, b []ScenarioSnapshot) (removed []ScenarioSnapshot) {
	for _, sa := range a {
		found := false
		for _, sb := range b {
			if reflect.DeepEqual(sa, sb) {
				found = true
				break
			}
		}
		if !found {
			removed = append(removed, sa)
		}
	}
	return
}
