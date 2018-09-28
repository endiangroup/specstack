package specstack

import (
	"github.com/endiangroup/specstack/actors"
	"github.com/endiangroup/specstack/repository"
)

type SpecStack interface {
	IsRepoInitialised() bool
}

func NewApp(repo repository.RepositoryReadWriter, developer actors.Developer) App {
	return App{
		Repo:      repo,
		Developer: developer,
	}
}

type App struct {
	Repo      repository.RepositoryReadWriter
	Developer actors.Developer
}

func (a App) IsRepoInitialised() bool {
	return a.Repo.IsInitialised()
}
