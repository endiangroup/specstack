package specstack

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/errors"
	"github.com/endiangroup/specstack/metadata"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/personas"
	"github.com/endiangroup/specstack/repository"
	"github.com/endiangroup/specstack/specification"
	"github.com/spf13/afero"
)

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
	GetStoryMetadata(string) ([]*metadata.Entry, error)
	RunRepoPostCommitHook() error
	RunRepoPostUpdateHook() error
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
	if a.isFirstRun(err) {
		return a.createDefaultConfig()
	} else if err != nil {
		return nil, err
	}

	return c, nil
}

func (a *appController) isFirstRun(err error) bool {
	return err == persistence.ErrNoConfigFound
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
		if err == persistence.ErrNoConfigFound {
			return MissingRequiredConfigValueErr(userName)
		}

		return err
	}

	c.User.Email, err = a.repo.GetConfig(userEmail)
	if err != nil {
		if err == persistence.ErrNoConfigFound {
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

func (a *appController) findStoryObject(name string) (*specification.Story, io.Reader, error) {
	reader := a.specificationReader()

	spec, warnings, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	for _, warning := range warnings {
		a.emitWarning(warning)
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

func (a *appController) GetStoryMetadata(storyName string) ([]*metadata.Entry, error) {
	_, object, err := a.findStoryObject(storyName)
	if err != nil {
		return nil, err
	}

	return metadata.ReadAll(a.omniStore, object)
}

func (a *appController) RunRepoPostCommitHook() error {
	if a.config.Project.PushingMode != config.ModeSemiAuto {
		return nil
	}
	return a.Push()
}

// FIXME! Rename to merge
func (a *appController) RunRepoPostUpdateHook() error {
	if a.config.Project.PullingMode != config.ModeSemiAuto {
		return nil
	}
	return a.Pull()
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
