package specstack

import (
	"context"
	"fmt"
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
}

func New(
	path string,
	repo repository.Repository,
	developer personas.Developer,
	configStore config.Storer,
	metadataStore metadata.ReadStorer,
) Controller {
	return &appController{
		path:          path,
		repo:          repo,
		developer:     developer,
		configStore:   configStore,
		metadataStore: metadataStore,
	}
}

type appController struct {
	path          string
	repo          repository.Repository
	configStore   config.Storer
	developer     personas.Developer
	config        *config.Config
	metadataStore metadata.ReadStorer
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
	c, err := config.Load(a.configStore)
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

	return config.Create(a.configStore, c)
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

// TODO: Move this to developer persona. How best to transfer deps?
func (a *appController) AddMetadataToStory(storyName, key, value string) error {

	reader := a.specificationReader()
	spec, warnings, err := reader.Read()

	if err != nil {
		return err
	}

	// FIXME! Emit warnings properly
	for _, warning := range warnings {
		fmt.Printf("WARNING: %s\n", warning.Error())
	}

	story, err := spec.FindStory(storyName)
	if err != nil {
		return err
	}

	object, err := reader.ReadSource(story)
	if err != nil {
		return err
	}

	entry := &metadata.Entry{
		Name:  key,
		Value: value,
	}

	return a.metadataStore.Store(object, entry)
}
