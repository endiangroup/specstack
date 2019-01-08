package specstack

import (
	"errors"

	"github.com/endiangroup/specstack/metadata"
	"github.com/endiangroup/specstack/repository"
)

var (
	// Thrown when a git repo is not initialised
	ErrUninitialisedRepo = errors.New("Please initialise repository first before running")
)

type ConfigAsserter interface {
	AssertConfig() error
}

type ConfigGetListSetter interface {
	ListConfiguration() (map[string]string, error)
	GetConfiguration(name string) (string, error)
	SetConfiguration(name, value string) error
}

type MetadataGetAdder interface {
	AddMetadataToStory(storyName, key, value string) error
	AddMetadataToScenario(scenarioName, storyName, key, value string) error
	GetStoryMetadata(string) ([]*metadata.Entry, error)
	GetScenarioMetadata(scenario string, story string) ([]*metadata.Entry, error)
}

type PushPuller interface {
	Push() error
	Pull() error
}

type MetadataTransferer interface {
	TransferScenarioMetadata() error
}

type RepoHooker interface {
	RepoPrePushHook() error
	RepoPostMergeHook() error
	RepoPostCommitHook() error
}

type Application struct {
	ConfigAsserter      ConfigAsserter
	ConfigGetListSetter ConfigGetListSetter
	Repository          repository.Repository
	MetadataGetAdder    MetadataGetAdder
	PushPuller          PushPuller
	MetadataTransferer  MetadataTransferer
	RepoHooker          RepoHooker
}

func (a *Application) Initialise() error {
	if !a.Repository.IsInitialised() {
		return ErrUninitialisedRepo
	}

	if err := a.ConfigAsserter.AssertConfig(); err != nil {
		return err
	}

	return metadata.PrepareSync(a.Repository)
}
