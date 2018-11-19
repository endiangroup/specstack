package main

import (
	"os"

	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/cmd"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/personas"
	"github.com/endiangroup/specstack/repository"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gitRepo := repository.NewGitRepository(dir)
	repoStore := persistence.NewStore(
		persistence.NewNamespacedKeyValueStorer(gitRepo, "specstack"),
		gitRepo,
	)
	developer := personas.NewDeveloper(repoStore, gitRepo)
	app := specstack.New(
		dir,
		gitRepo,
		developer,
		repoStore,
	)
	cobra := cmd.WireUpCobraHarness(
		cmd.NewCobraHarness(app, os.Stdin, os.Stdout, os.Stderr),
	)

	if err := cobra.Execute(); err != nil {
		if cliErr, ok := err.(cmd.CliErr); ok {
			os.Exit(cliErr.ExitCode)
		}

		os.Exit(-1)
	}

	os.Exit(0)
}
