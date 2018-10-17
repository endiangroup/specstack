package main

import (
	"os"

	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/actors"
	"github.com/endiangroup/specstack/cmd"
	"github.com/endiangroup/specstack/persistence"
	"github.com/endiangroup/specstack/repository"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gitRepo := repository.NewGitRepo(dir, "specstack")
	repoStore := persistence.NewRepositoryStore(gitRepo)
	developer := actors.NewDeveloper(repoStore)
	app := specstack.NewApp(dir, gitRepo, developer, repoStore)

	cobra := cmd.WireUpHarness(cmd.NewCobraHarness(app, os.Stdin, os.Stdout, os.Stderr))

	if err := cobra.Execute(); err != nil {
		if cliErr, ok := err.(cmd.CliErr); ok {
			os.Exit(cliErr.ExitCode)
		}

		os.Exit(-1)
	}

	os.Exit(0)
}
