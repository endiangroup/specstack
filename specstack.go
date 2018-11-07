package specstack

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/personas"
	"github.com/endiangroup/specstack/repository"
)

var (
	// Thrown when a git repo is not initialised
	ErrUninitialisedRepo = errors.New("Please initialise repository first before running")
)

type Controller interface {
	Initialise() error
	ListConfiguration() (map[string]string, error)
	GetConfiguration(string) (string, error)
	SetConfiguration(string, string) error
}

func New(path string, repo repository.Initialiser, developer personas.Developer, configStore config.Storer) Controller {
	return &appController{
		path:        path,
		repo:        repo,
		developer:   developer,
		configStore: configStore,
	}
}

type appController struct {
	path        string
	repo        repository.Initialiser
	configStore config.Storer
	developer   personas.Developer
	config      *config.Config
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
	c.Project.Name = filepath.Base(a.path)

	return config.Create(a.configStore, c)
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
