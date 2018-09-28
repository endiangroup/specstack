package main

import (
	"os"

	"github.com/endiangroup/specstack"
	"github.com/endiangroup/specstack/actors"
	"github.com/endiangroup/specstack/cmd"
	"github.com/endiangroup/specstack/config"
	"github.com/endiangroup/specstack/repository"
)

func main() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	gitRepo := repository.NewGitRepo(dir)
	config := config.NewRepositoryConfig(gitRepo)
	developer := actors.NewDeveloper(config)
	app := specstack.NewApp(gitRepo, developer)

	cmd.WireUpHarness(cmd.NewCobraHarness(app, os.Stdin, os.Stdout, os.Stderr))

	if err := cmd.Root.Execute(); err != nil {
		if cliErr, ok := err.(cmd.CliErr); ok {
			os.Exit(cliErr.ReturnCode)
		}

		os.Exit(-1)
	}

	os.Exit(0)
}
