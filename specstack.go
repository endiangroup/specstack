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

	return err
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

func (a *appController) warning(warning error) {
	fmt.Fprintf(a.stderr, "WARNING: %s\n", warning.Error())
}

func (a *appController) findStoryObject(name string) (*specification.Story, io.Reader, error) {
	reader := a.specificationReader()

	spec, warnings, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}

	for _, warning := range warnings {
		a.warning(warning)
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

	return a.developer.AddMetadataToStory(a.newContextWithConfig(), story, object, key, value)
}

func (a *appController) GetStoryMetadata(storyName string) ([]*metadata.Entry, error) {
	_, object, err := a.findStoryObject(storyName)
	if err != nil {
		return nil, err
	}

	return metadata.ReadAll(a.omniStore, object)
}
