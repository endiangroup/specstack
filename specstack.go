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
	ErrUninitialisedRepo = errors.New("Please initialise repository first before running")
)

type SpecStack interface {
	Initialise() error
	ListConfiguration() (map[string]string, error)
	GetConfiguration(string) (string, error)
	SetConfiguration(string, string) error
}

func NewApp(path string, repo repository.Initialiser, developer personas.Developer, configStore config.Storer) *App {
	return &App{
		path:        path,
		repo:        repo,
		developer:   developer,
		configStore: configStore,
	}
}

type App struct {
	path        string
	repo        repository.Initialiser
	configStore config.Storer
	developer   personas.Developer
	config      *config.Config
}

func (a *App) Initialise() error {
	if !a.repo.IsInitialised() {
		return ErrUninitialisedRepo
	}

	var err error
	if a.config, err = a.loadOrCreateConfig(); err != nil {
		return err
	}

	return nil
}

func (a *App) loadOrCreateConfig() (*config.Config, error) {
	c, err := config.Load(a.configStore)
	if a.isFirstRun(err) {
		return a.createDefaultConfig()
	} else if err != nil {
		return nil, err
	}

	return c, nil
}

func (a *App) isFirstRun(err error) bool {
	return err == persistence.ErrNoConfigFound
}

func (a *App) createDefaultConfig() (*config.Config, error) {
	c := config.NewWithDefaults()
	c.Project.Name = filepath.Base(a.path)

	return config.Create(a.configStore, c)
}

func (a *App) ListConfiguration() (map[string]string, error) {
	return a.developer.ListConfiguration(a.newContextWithConfig())
}

func (a *App) GetConfiguration(name string) (string, error) {
	return a.developer.GetConfiguration(a.newContextWithConfig(), name)
}

func (a *App) SetConfiguration(name, value string) error {
	return a.developer.SetConfiguration(a.newContextWithConfig(), name, value)
}

func (a *App) newContextWithConfig() context.Context {
	return config.InContext(context.TODO(), a.config)
}
