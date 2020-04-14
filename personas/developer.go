package personas

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/errors"
	"github.com/endiangroup/specstack/metadata"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/repository"
	"github.com/endiangroup/specstack/snapshot"
	"github.com/endiangroup/specstack/specification"
	"github.com/spf13/afero"
)

type MissingRequiredConfigValueErr string

func (err MissingRequiredConfigValueErr) Error() string {
	return fmt.Sprintf("no %s set in repository", string(err))
}

type Developer struct {
	path   string
	store  *persistence.Store
	config *config.Config
	repo   repository.Repository
	stdout io.Writer
	stderr io.Writer
}

func NewDeveloper(
	path string,
	store *persistence.Store,
	repo repository.Repository,
	stdout, stderr io.Writer,
) *Developer {
	return &Developer{
		path:   path,
		store:  store,
		repo:   repo,
		stdout: stdout,
		stderr: stderr,
	}
}

func (d *Developer) AssertConfig() error {
	c, err := d.loadOrCreateConfig()
	if err != nil {
		return err
	}
	d.config = c
	return nil
}

func (d *Developer) loadOrCreateConfig() (*config.Config, error) {
	c, err := config.Load(d.store)
	if d.isErrConfigNotFound(err) {
		return d.createDefaultConfig()
	} else if err != nil {
		return nil, err
	}

	return c, nil
}

func (d *Developer) isErrConfigNotFound(err error) bool {
	return err == persistence.ErrNoConfigFound || err == repository.ErrNoConfigFound
}

func (d *Developer) createDefaultConfig() (*config.Config, error) {
	c := config.NewWithDefaults()

	d.setProjectDefaults(c)
	if err := d.setUserDefaults(c); err != nil {
		return nil, err
	}

	return config.Create(d.store, c)
}

func (d *Developer) setProjectDefaults(c *config.Config) {
	c.Project.Name = filepath.Base(d.path)
}

func (d *Developer) setUserDefaults(c *config.Config) error {
	var err error
	userName := config.KeyUser.Append(config.KeyUserName)
	userEmail := config.KeyUser.Append(config.KeyUserEmail)

	c.User.Name, err = d.repo.GetConfig(userName)
	if err != nil {
		if d.isErrConfigNotFound(err) {
			return MissingRequiredConfigValueErr(userName)
		}

		return err
	}

	c.User.Email, err = d.repo.GetConfig(userEmail)
	if err != nil {
		if d.isErrConfigNotFound(err) {
			return MissingRequiredConfigValueErr(userEmail)
		}

		return err
	}

	return nil
}

func (d *Developer) ListConfiguration() (map[string]string, error) {
	return config.ToMap(d.config), nil
}

func (d *Developer) GetConfiguration(name string) (string, error) {
	return config.Get(d.config, name)
}

func (d *Developer) SetConfiguration(name, value string) error {
	err := config.Set(d.config, name, value)
	if err != nil {
		return err
	}

	_, err = config.Store(d.store, d.config)
	return err
}

func (d *Developer) specificationFactory() *specification.Factory {
	return specification.NewFactory(
		afero.NewOsFs(),
		d.config.Project.FeaturesDir,
		d.stderr,
	)
}

func (d *Developer) specification() (*specification.Specification, specification.Reader, error) {
	return d.specificationFactory().Specification()
}

func (d *Developer) findStoryObject(name string) (*specification.Story, io.Reader, error) {
	spec, reader, err := d.specification()
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

func (d *Developer) findScenarioObject(name, story string) (*specification.Scenario, io.Reader, error) {
	spec, reader, err := d.specification()
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

func (d *Developer) AddMetadataToStory(storyName, key, value string) error {
	_, object, err := d.findStoryObject(storyName)
	if err != nil {
		return err
	}

	if err := metadata.Add(d.store, object, metadata.NewKeyValue(key, value)); err != nil {
		return err
	}

	if d.config.Project.PushingMode == config.ModeAuto {
		return errors.WarningOrNil(d.Push())
	}

	return nil
}

func (d *Developer) AddMetadataToScenario(name, storyName, key, value string) error {
	_, object, err := d.findScenarioObject(name, storyName)
	if err != nil {
		return err
	}

	if err := metadata.Add(d.store, object, metadata.NewKeyValue(key, value)); err != nil {
		return err
	}

	if d.config.Project.PushingMode == config.ModeAuto {
		return errors.WarningOrNil(d.Push())
	}

	return nil
}

func (d *Developer) GetStoryMetadata(name string) ([]*metadata.Entry, error) {
	_, object, err := d.findStoryObject(name)
	if err != nil {
		return nil, err
	}

	return metadata.ReadAll(d.store, object)
}

func (d *Developer) GetScenarioMetadata(name, story string) ([]*metadata.Entry, error) {
	_, object, err := d.findScenarioObject(name, story)
	if err != nil {
		return nil, err
	}

	return metadata.ReadAll(d.store, object)
}

func (d *Developer) Pull() error {
	if d.config.Project.Remote == "" {
		return fmt.Errorf("configure a project remote first")
	}
	return metadata.Pull(d.repo, d.config.Project.Remote)
}

func (d *Developer) Push() error {
	if d.config.Project.Remote == "" {
		return fmt.Errorf("configure a project remote first")
	}
	return metadata.Push(d.repo, d.config.Project.Remote)
}

func (d *Developer) TransferScenarioMetadata() error {
	ss := snapshot.NewScenarioMetadataSnapshotter(
		d.specificationFactory(),
		d.store,
		"snapshots",
		d.repo,
		d.config.Project.FeaturesDir,
	)
	return ss.Snapshot()
}

func (d *Developer) RepoPrePushHook() error {
	if d.config.Project.PushingMode != config.ModeSemiAuto {
		return nil
	}
	if err := d.TransferScenarioMetadata(); err != nil {
		return err
	}
	return d.Push()
}

func (d *Developer) RepoPostMergeHook() error {
	if d.config.Project.PullingMode != config.ModeSemiAuto {
		return nil
	}
	if err := d.Pull(); err != nil {
		return err
	}
	return d.TransferScenarioMetadata()
}

func (d *Developer) RepoPostCommitHook() error {
	return d.TransferScenarioMetadata()
}

func (d *Developer) RepoRemovePrePushHook() error {
	return d.repo.RemoveHook("pre-push")
}
func (d *Developer) RepoRemovePostMergeHook() error {
	return d.repo.RemoveHook("post-merge")
}
func (d *Developer) RepoRemovePostCommitHook() error {
	return d.repo.RemoveHook("post-commit")
}
