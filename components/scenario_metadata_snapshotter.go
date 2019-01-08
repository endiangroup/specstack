package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/endiangroup/specstack/fuzzy"
	"github.com/endiangroup/specstack/metadata"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/repository"
	"github.com/endiangroup/specstack/specification"
	"github.com/spf13/afero"
)

type ScenarioMetadataSnapshotter struct {
	Factory     *SpecificationFactory
	Store       *persistence.Store
	StorageKey  string
	Repository  repository.Repository
	FeaturesDir string
}

func NewScenarioMetadataSnapshotter(
	factory *SpecificationFactory,
	store *persistence.Store,
	storageKey string,
	repository repository.Repository,
	featuresDir string,
) *ScenarioMetadataSnapshotter {
	return &ScenarioMetadataSnapshotter{
		Factory:     factory,
		Store:       store,
		StorageKey:  storageKey,
		Repository:  repository,
		FeaturesDir: featuresDir,
	}
}

func (s *ScenarioMetadataSnapshotter) storageKeyReader() io.Reader {
	return bytes.NewBufferString(s.StorageKey)
}

func (s *ScenarioMetadataSnapshotter) scenarioHasMetadata(scenario *specification.Scenario) bool {
	reader := s.Factory.SpecificationReader()
	key, err := reader.ReadSource(scenario)
	if err != nil {
		return false
	}
	e, err := metadata.ReadAll(s.Store, key)
	return err == nil && len(e) > 0
}

func (s *ScenarioMetadataSnapshotter) previousSnapshot() (specification.Snapshot, error) {
	entries, err := metadata.ReadAll(s.Store, s.storageKeyReader())
	if err != nil || len(entries) == 0 {
		return specification.Snapshot{}, err
	}

	latest := entries[len(entries)-1]
	snap := specification.Snapshot{}
	if err := json.Unmarshal([]byte(latest.Value), &snap); err != nil {
		return snap, err
	}

	return snap, nil
}

func (s *ScenarioMetadataSnapshotter) currentSnapshot() (specification.Snapshot, error) {
	spec, reader, err := s.Factory.Specification()
	if err != nil {
		return specification.Snapshot{}, err
	}

	snapshotter := specification.NewSnapshotter(reader, s.Repository)
	current, err := snapshotter.Snapshot(spec)
	if err != nil {
		return specification.Snapshot{}, err
	}
	return current, nil
}

func (s *ScenarioMetadataSnapshotter) snapshots() (current, previous specification.Snapshot, err error) {
	current, err = s.currentSnapshot()
	if err != nil {
		return
	}
	previous, err = s.previousSnapshot()
	return
}

func (s *ScenarioMetadataSnapshotter) storeSnapshot(snap specification.Snapshot) error {
	jsn, err := json.Marshal(snap)
	if err != nil {
		return err
	}
	return metadata.Add(s.Store, s.storageKeyReader(), metadata.NewKeyValue("snapshot", string(jsn)))
}

/*
Loads a scenario from a snapshot. The procedure is:

1. Look in git object for story file
2. If not there, look on disk
3. If found, create new Scenario with only that one story file
4. Query spec for scenario at line number
5. Return if found
*/
func (s *ScenarioMetadataSnapshotter) scenarioFromSnapshot(snap specification.ScenarioSnapshot) (*specification.Scenario, error) {
	fs, err := s.fileSystemFromScenarioSnapshot(snap)
	if err != nil {
		return nil, err
	}

	reader := specification.NewFilesystemReader(fs, s.FeaturesDir)
	spec, _, err := reader.Read()
	if err != nil {
		return nil, err
	}

	scenarios := specification.NewQuery(spec).MapReduce(
		specification.MapScenarios(),
		specification.MapScenarioLineNumber(snap.LineNumber),
	).Scenarios()

	if l := len(scenarios); l != 1 {
		return nil, fmt.Errorf("Expected 1 scenario from query, got %d", l)
	}

	return scenarios[0], nil
}

func (s *ScenarioMetadataSnapshotter) scenarioMapFromSnapshots(
	reader specification.Reader,
	snapshots []specification.ScenarioSnapshot,
) (map[io.Reader]*specification.Scenario, error) {
	output := make(map[io.Reader]*specification.Scenario)
	for _, sn := range snapshots {
		scen, err := s.scenarioFromSnapshot(sn)
		if err != nil {
			continue
		}
		if !s.scenarioHasMetadata(scen) {
			continue
		}
		key, err := reader.ReadSource(scen)
		if err != nil {
			return nil, err
		}

		output[key] = scen
	}
	return output, nil
}

func (s *ScenarioMetadataSnapshotter) fileSystemFromScenarioSnapshot(snap specification.ScenarioSnapshot) (afero.Fs, error) {
	fs := afero.NewMemMapFs()
	var (
		err         error
		fileContent string
	)

	fileContent, err = s.Repository.ObjectString(snap.StoryID)

	if err != nil && snap.StorySource.Type == specification.SourceTypeFile {
		if fc, err := ioutil.ReadFile(snap.StorySource.Body); err == nil {
			fileContent = string(fc)
		}
	}

	if fileContent != "" {
		err := afero.WriteFile(fs, snap.StorySource.Body, []byte(fileContent), os.ModePerm)
		return fs, err
	}

	return nil, fmt.Errorf("spec not found")
}

func (s *ScenarioMetadataSnapshotter) scenarioParent(
	to *specification.Scenario,
	from map[io.Reader]*specification.Scenario,
) (*specification.Scenario, io.Reader) {
	var (
		bestDistance float64
		bestParent   *specification.Scenario
		parentObject io.Reader
	)
	for k, v := range from {
		if distance := specification.ScenarioDistance(to, v); distance >= fuzzy.DistanceThreshold &&
			distance > bestDistance {
			bestDistance = distance
			bestParent = v
			parentObject = k
		}
	}
	return bestParent, parentObject
}

func (s *ScenarioMetadataSnapshotter) transferScenarioMetadata(
	reader specification.Reader,
	scenario *specification.Scenario,
	potentialParents map[io.Reader]*specification.Scenario) error {
	if bestParent, parentObject := s.scenarioParent(scenario, potentialParents); bestParent != nil {
		object, err := reader.ReadSource(scenario)
		if err != nil {
			return err
		}

		if err := s.transferMetadata(parentObject, object); err != nil {
			return err
		}
	}
	return nil
}

func (s *ScenarioMetadataSnapshotter) transferMetadata(fromObject, toObject io.Reader) error {
	entries, err := metadata.ReadAll(s.Store, fromObject)
	if err != nil {
		return err
	}
	return metadata.Add(s.Store, toObject, entries...)
}

func (s *ScenarioMetadataSnapshotter) Snapshot() error {
	current, previous, err := s.snapshots()
	if err != nil {
		return err
	}

	if previous.Equal(current) {
		return nil
	}

	if err := s.storeSnapshot(current); err != nil {
		return err
	}

	removed, added := previous.Diff(current)
	if len(added.Scenarios) == 0 || len(removed.Scenarios) == 0 {
		return nil
	}

	reader := s.Factory.SpecificationReader()
	removedScenarios, err := s.scenarioMapFromSnapshots(reader, removed.Scenarios)
	if err != nil {
		return err
	}

	for _, snapshot := range added.Scenarios {
		scenario, err := s.scenarioFromSnapshot(snapshot)
		if err != nil {
			return err
		}
		if err := s.transferScenarioMetadata(reader, scenario, removedScenarios); err != nil {
			return err
		}
	}
	return nil
}
