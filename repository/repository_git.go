package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type GitRepo struct {
	Path string
}

func NewGitRepo(path string) *GitRepo {
	return &GitRepo{Path: path}
}

func (repo *GitRepo) IsInitialised() bool {
	_, _, err := repo.runGitCommandRaw("rev-parse")

	return err == nil
}

func (repo *GitRepo) Init() error {
	_, err := repo.runGitCommand("init")
	return err
}

func (repo *GitRepo) ConfigGetRegex(regex string) (string, error) {
	return repo.runGitCommand("config", "--get-regex", regex)
}

func (repo *GitRepo) runGitCommandRaw(args ...string) (string, string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.Path
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), err
}

func (repo *GitRepo) runGitCommand(args ...string) (string, error) {
	stdout, stderr, err := repo.runGitCommandRaw(args...)
	if err != nil {
		if stderr == "" {
			stderr = "Error running git command: " + strings.Join(args, " ")
		}
		err = fmt.Errorf(stderr)
	}
	return stdout, err
}
