package specstack

import (
	"errors"
	"path/filepath"

	"github.com/endiangroup/specstack/actors"
	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/repository"
)

var (
	ErrUninitialisedRepo = errors.New("Please initialise repository first before running")
)

type SpecStack interface {
	Initialise() error
	ListConfiguration() (map[string]string, error)
}

func NewApp(path string, repo repository.ReadWriter, developer actors.Developer, configStore config.Storer) App {
	return App{
		path:        path,
		repo:        repo,
		developer:   developer,
		configStore: configStore,
	}
}

type App struct {
	path        string
	repo        repository.ReadWriter
	configStore config.Storer
	developer   actors.Developer
}

func (a App) Initialise() error {
	if !a.repo.IsInitialised() {
		return ErrUninitialisedRepo
	}

	if a.isFirstRun() {
		if err := a.createDefaultConfig(); err != nil {
			return err
		}
	}

	return nil
}

func (a App) createDefaultConfig() error {
	c := config.NewWithDefaults()
	c.Project.Name = filepath.Base(a.path)

	_, err := a.configStore.CreateConfig(c)

	return err
}

func (a App) isFirstRun() bool {
	if _, err := a.ListConfiguration(); err != nil {
		if err == persistence.ErrNoConfigFound {
			return true
		}
	}

	return false
}

func (a App) ListConfiguration() (map[string]string, error) {
	return a.developer.ListConfiguration()
}
