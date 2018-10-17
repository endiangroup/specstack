package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func NewGitConfigErr(gitCmdErr *GitCmdErr) error {
	switch gitCmdErr.ExitCode {
	case 1:
		return GitConfigMissingSectionKeyErr{gitCmdErr}
	}

	return gitCmdErr
}

type GitConfigMissingSectionKeyErr struct {
	*GitCmdErr
}

func NewGitCmdErr(stderr string, exitCode int, args ...string) *GitCmdErr {
	return &GitCmdErr{
		Stderr:   stderr,
		ExitCode: exitCode,
		Args:     args,
	}
}

type GitCmdErr struct {
	Stderr   string
	ExitCode int
	Args     []string
}

func (err GitCmdErr) Error() string {
	if err.Stderr == "" {
		return fmt.Sprintf("error running git command (exit code: %d): %s", err.ExitCode, strings.Join(err.Args, " "))
	}

	return err.Stderr
}

type GitRepo struct {
	Path            string
	ConfigNamespace string
}

func NewGitRepo(path, configNamespace string) *GitRepo {
	return &GitRepo{
		Path:            path,
		ConfigNamespace: configNamespace,
	}
}

func (repo *GitRepo) IsInitialised() bool {
	_, _, _, err := repo.runGitCommandRaw("rev-parse")

	return err == nil
}

func (repo *GitRepo) Init() error {
	_, err := repo.runGitCommand("init")
	return err
}

func (repo *GitRepo) ConfigGetAll() (map[string]string, error) {
	result, err := repo.runGitCommand("config", "--null", "--get-regex", "^"+repo.ConfigNamespace+`\.`)
	if err != nil {
		return nil, NewGitConfigErr(err.(*GitCmdErr))
	}

	configMap := map[string]string{}
	for _, kvPair := range strings.Split(result, "\x00") {
		kvPair = strings.TrimSpace(kvPair)

		if kvPair == "" {
			continue
		}
		kvParts := strings.SplitN(kvPair, "\n", 2)
		if len(kvParts) == 1 {
			configMap[repo.trimConfigNamespace(kvParts[0])] = ""
		} else {
			configMap[repo.trimConfigNamespace(kvParts[0])] = kvParts[1]
		}
	}

	return configMap, nil
}

func (repo *GitRepo) ConfigGet(key string) (string, error) {
	return repo.runGitCommand("config", "--get", repo.prefixConfigNamespace(key))
}

func (repo *GitRepo) ConfigSet(key, value string) error {
	_, err := repo.runGitCommand("config", repo.prefixConfigNamespace(key), value)
	return err
}
func (repo *GitRepo) ConfigUnset(key string) error {
	_, err := repo.runGitCommand("config", "--unset", repo.prefixConfigNamespace(key))

	return err
}

func (repo *GitRepo) runGitCommandRaw(args ...string) (string, string, int, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.Path
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	var exitCode int
	if err != nil {
		exitCode = 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
	}

	return strings.TrimSpace(stdout.String()), strings.TrimSpace(stderr.String()), exitCode, err
}

func (repo *GitRepo) runGitCommand(args ...string) (string, error) {
	stdout, stderr, exitCode, err := repo.runGitCommandRaw(args...)
	if err != nil {
		return "", NewGitCmdErr(stderr, exitCode, args...)
	}

	return stdout, err
}

func (repo *GitRepo) prefixConfigNamespace(key string) string {
	return repo.ConfigNamespace + "." + key
}

func (repo *GitRepo) trimConfigNamespace(key string) string {
	return strings.TrimPrefix(key, repo.ConfigNamespace+".")
}
