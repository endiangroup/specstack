package specstack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/errors"
	"github.com/endiangroup/specstack/fuzzy"
	"github.com/endiangroup/specstack/metadata"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/personas"
	"github.com/endiangroup/specstack/repository"
	"github.com/endiangroup/specstack/specification"
	"github.com/spf13/afero"
)

func MetdataSnapshotKey() io.Reader {
	return bytes.NewBufferString("snapshots")
}

type MissingRequiredConfigValueErr string

func (err MissingRequiredConfigValueErr) Error() string {
	return fmt.Sprintf("no %s set in repository", string(err))
}

var (
	// Thrown when a git repo is not initialised
	ErrUninitialisedRepo = errors.New("Please initialise repository first before running")
)

type Controller interface {
	Initialise() error
	ListConfiguration() (map[string]string, error)
	GetConfiguration(string) (string, error)
	SetConfiguration(string, string) error
	AddMetadataToStory(storyName, key, value string) error
	AddMetadataToScenario(scenarioName, storyName, key, value string) error
	GetStoryMetadata(string) ([]*metadata.Entry, error)
	GetScenarioMetadata(scenario string, story string) ([]*metadata.Entry, error)
	SnapshotScenarioMetadata() error
	RunRepoPrePushHook() error
	RunRepoPostMergeHook() error
	RunRepoPostCommitHook() error
	Push() error
	Pull() error
}

func New(
	path string,
	repo repository.Repository,
	developer personas.Developer,
	omniStore *persistence.Store,
	stdout io.Writer,
	stderr io.Writer,
) Controller {
	return &appController{
		path:      path,
		repo:      repo,
		developer: developer,
		omniStore: omniStore,
		stdout:    stdout,
		stderr:    stderr,
	}
}

type appController struct {
	path      string
	repo      repository.Repository
	omniStore *persistence.Store
	developer personas.Developer
	config    *config.Config
	stdout    io.Writer
	stderr    io.Writer
}

func (a *appController) Initialise() error {
	if !a.repo.IsInitialised() {
		return ErrUninitialisedRepo
	}

	var err error
	a.config, err = a.loadOrCreateConfig()
	if err != nil {
		return err
	}

	return metadata.PrepareSync(a.repo)
}

func (a *appController) loadOrCreateConfig() (*config.Config, error) {
	c, err := config.Load(a.omniStore)
	if a.isErrConfigNotFound(err) {
		return a.createDefaultConfig()
	} else if err != nil {
		return nil, err
	}

	return c, nil
}

func (a *appController) isErrConfigNotFound(err error) bool {
	return err == persistence.ErrNoConfigFound || err == repository.ErrNoConfigFound
}

func (a *appController) createDefaultConfig() (*config.Config, error) {
	c := config.NewWithDefaults()

	a.setProjectDefaults(c)
	if err := a.setUserDefaults(c); err != nil {
		return nil, err
	}

	return config.Create(a.omniStore, c)
}

func (a *appController) setProjectDefaults(c *config.Config) {
	c.Project.Name = filepath.Base(a.path)
}

func (a *appController) setUserDefaults(c *config.Config) error {
	var err error
	userName := config.KeyUser.Append(config.KeyUserName)
	userEmail := config.KeyUser.Append(config.KeyUserEmail)

	c.User.Name, err = a.repo.GetConfig(userName)
	if err != nil {
		if a.isErrConfigNotFound(err) {
			return MissingRequiredConfigValueErr(userName)
		}

		return err
	}

	c.User.Email, err = a.repo.GetConfig(userEmail)
	if err != nil {
		if a.isErrConfigNotFound(err) {
			return MissingRequiredConfigValueErr(userEmail)
		}

		return err
	}

	return nil
}

func (a *appController) ListConfiguration() (map[string]string, error) {
	return a.developer.ListConfiguration(a.newContextWithConfig())
}

func (a *appController) GetConfiguration(name string) (string, error) {
	return a.developer.GetConfiguration(a.newContextWithConfig(), name)
}

func (a *appController) SetConfiguration(name, value string) error {
	return a.developer.SetConfiguration(a.newContextWithConfig(), name, value)
}

func (a *appController) newContextWithConfig() context.Context {
	return config.InContext(context.TODO(), a.config)
}

func (a *appController) specificationReader() specification.Reader {
	return specification.NewFilesystemReader(afero.NewOsFs(), a.config.Project.FeaturesDir)
}

func (a *appController) emitWarning(warning error) {
	fmt.Fprintf(a.stderr, "WARNING: %s\n", warning.Error())
}

func (a *appController) specification() (*specification.Specification, specification.Reader, error) {
	reader := a.specificationReader()

	spec, warnings, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	for _, warning := range warnings {
		a.emitWarning(warning)
	}
	return spec, reader, nil
}

func (a *appController) findStoryObject(name string) (*specification.Story, io.Reader, error) {
	spec, reader, err := a.specification()
	if err != nil {
		return nil, nil, err
	}

	story, err := spec.FindStory(name)
	if err != nil {
		return nil, nil, err
	}

	object, err := reader.ReadSource(story)
	if err != nil {
		return nil, nil, err
	}

	return story, object, nil
}

func (a *appController) findScenarioObject(name, story string) (*specification.Scenario, io.Reader, error) {
	spec, reader, err := a.specification()
	if err != nil {
		return nil, nil, err
	}
	scenario, err := spec.FindScenario(name, story)
	if err != nil {
		return nil, nil, err
	}

	object, err := reader.ReadSource(scenario)
	if err != nil {
		return nil, nil, err
	}

	return scenario, object, nil
}

func (a *appController) AddMetadataToStory(storyName, key, value string) error {
	story, object, err := a.findStoryObject(storyName)
	if err != nil {
		return err
	}

	if err := a.developer.AddMetadataToStory(
		a.newContextWithConfig(),
		story,
		object,
		key,
		value,
	); err != nil {
		return err
	}

	if a.config.Project.PushingMode == config.ModeAuto {
		return errors.WarningOrNil(a.Push())
	}

	return nil
}

func (a *appController) AddMetadataToScenario(name, storyName, key, value string) error {
	scenario, object, err := a.findScenarioObject(name, storyName)
	if err != nil {
		return err
	}
	if err := a.developer.AddMetadataToScenario(
		a.newContextWithConfig(),
		scenario,
		object,
		key,
		value,
	); err != nil {
		return err
	}

	if a.config.Project.PushingMode == config.ModeAuto {
		return errors.WarningOrNil(a.Push())
	}

	return nil
}

func (a *appController) GetStoryMetadata(name string) ([]*metadata.Entry, error) {
	_, object, err := a.findStoryObject(name)
	if err != nil {
		return nil, err
	}

	return metadata.ReadAll(a.omniStore, object)
}

func (a *appController) GetScenarioMetadata(name, story string) ([]*metadata.Entry, error) {
	_, object, err := a.findScenarioObject(name, story)
	if err != nil {
		return nil, err
	}

	return metadata.ReadAll(a.omniStore, object)
}

func (a *appController) scenarioHasMetadata(scenario *specification.Scenario) bool {
	reader := a.specificationReader()
	key, err := reader.ReadSource(scenario)
	if err != nil {
		return false
	}
	e, err := metadata.ReadAll(a.omniStore, key)
	return err == nil && len(e) > 0
}

func (a *appController) previousSnapshot() (specification.Snapshot, error) {
	entries, err := metadata.ReadAll(a.omniStore, MetdataSnapshotKey())
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

func (a *appController) currentSnapshot() (specification.Snapshot, error) {
	spec, reader, err := a.specification()
	if err != nil {
		return specification.Snapshot{}, err
	}

	snapshotter := specification.NewSnapshotter(reader, a.repo)
	current, err := snapshotter.Snapshot(spec)
	if err != nil {
		return specification.Snapshot{}, err
	}
	return current, nil
}

func (a *appController) snapshots() (current, previous specification.Snapshot, err error) {
	current, err = a.currentSnapshot()
	if err != nil {
		return
	}
	previous, err = a.previousSnapshot()
	return
}

func (a *appController) storeSnapshot(s specification.Snapshot) error {
	jsn, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return metadata.Add(a.omniStore, MetdataSnapshotKey(), metadata.NewKeyValue("snapshot", string(jsn)))
}

/*
Loads a scenario from a snapshot. The procedure is:

1. Look in git object for story file
2. If not there, look on disk
3. If found, create new Scenario with only that one story file
4. Query spec for scenario at line number
5. Return if found
*/
func (a *appController) scenarioFromSnapshot(snap specification.ScenarioSnapshot) (*specification.Scenario, error) {
	fs, err := a.fileSystemFromScenarioSnapshot(snap)
	if err != nil {
		return nil, err
	}

	reader := specification.NewFilesystemReader(fs, a.config.Project.FeaturesDir)
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

func (a *appController) scenarioMapFromSnapshots(
	reader specification.Reader,
	snapshots []specification.ScenarioSnapshot,
) (map[io.Reader]*specification.Scenario, error) {
	output := make(map[io.Reader]*specification.Scenario)
	for _, s := range snapshots {
		scen, err := a.scenarioFromSnapshot(s)
		if err != nil {
			continue
		}
		if !a.scenarioHasMetadata(scen) {
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

func (a *appController) fileSystemFromScenarioSnapshot(snap specification.ScenarioSnapshot) (afero.Fs, error) {
	fs := afero.NewMemMapFs()
	var (
		err         error
		fileContent string
	)

	fileContent, err = a.repo.ObjectString(snap.StoryID)

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

func (a *appController) scenarioParent(
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

func (a *appController) transferScenarioMetadata(
	reader specification.Reader,
	scenario *specification.Scenario,
	potentialParents map[io.Reader]*specification.Scenario) error {
	if bestParent, parentObject := a.scenarioParent(scenario, potentialParents); bestParent != nil {
		object, err := reader.ReadSource(scenario)
		if err != nil {
			return err
		}

		if err := a.developer.TransferScenarioMetadata(
			bestParent, scenario,
			parentObject, object,
		); err != nil {
			return err
		}
	}
	return nil
}

func (a *appController) SnapshotScenarioMetadata() error {
	current, previous, err := a.snapshots()
	if err != nil {
		return err
	}

	if previous.Equal(current) {
		return nil
	}

	if err := a.storeSnapshot(current); err != nil {
		return err
	}

	removed, added := previous.Diff(current)
	if len(added.Scenarios) == 0 || len(removed.Scenarios) == 0 {
		return nil
	}

	reader := a.specificationReader()
	removedScenarios, err := a.scenarioMapFromSnapshots(reader, removed.Scenarios)
	if err != nil {
		return err
	}

	for _, snapshot := range added.Scenarios {
		scenario, err := a.scenarioFromSnapshot(snapshot)
		if err != nil {
			return err
		}
		if err := a.transferScenarioMetadata(reader, scenario, removedScenarios); err != nil {
			return err
		}
	}
	return nil
}

func (a *appController) RunRepoPrePushHook() error {
	if a.config.Project.PushingMode != config.ModeSemiAuto {
		return nil
	}
	if err := a.SnapshotScenarioMetadata(); err != nil {
		return err
	}
	return a.Push()
}

func (a *appController) RunRepoPostMergeHook() error {
	if a.config.Project.PullingMode != config.ModeSemiAuto {
		return nil
	}
	if err := a.Pull(); err != nil {
		return err
	}
	return a.SnapshotScenarioMetadata()
}

func (a *appController) RunRepoPostCommitHook() error {
	return a.SnapshotScenarioMetadata()
}

func (a *appController) Pull() error {
	if a.config.Project.Remote == "" {
		return fmt.Errorf("configure a project remote first")
	}
	return metadata.Pull(a.repo, a.config.Project.Remote)
}

func (a *appController) Push() error {
	if a.config.Project.Remote == "" {
		return fmt.Errorf("configure a project remote first")
	}
	return metadata.Push(a.repo, a.config.Project.Remote)
}
