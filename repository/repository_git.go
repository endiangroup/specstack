package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

const (
	gitNotesRef = "refs/notes/specstack"

	// Scopes for git config
	GitConfigScopeLocal  = 1
	GitConfigScopeSystem = 2
	GitConfigScopeGlobal = 4
)

// NewGitConfigErr creates the appropriate typed error for a Git failure, if
// possible.
func NewGitConfigErr(gitCmdErr *GitCmdErr) error {
	switch gitCmdErr.ExitCode {
	case 1:
		return GitConfigMissingKeyErr{gitCmdErr}
	}

	return gitCmdErr
}

type GitConfigMissingKeyErr struct {
	*GitCmdErr
}

// NewGitCmdErr creates an instance of an error used for failed Git commands
func NewGitCmdErr(stderr string, exitCode int, args ...string) error {
	return &GitCmdErr{
		Stderr:   stderr,
		ExitCode: exitCode,
		Args:     args,
	}
}

// GitCmdErr is an error types for failed Git commands.
type GitCmdErr struct {
	Stderr   string
	ExitCode int
	Args     []string
}

func (err GitCmdErr) Error() string {
	if err.Stderr == "" {
		return fmt.Sprintf(
			"error running git command (exit code: %d): %s",
			err.ExitCode,
			strings.Join(err.Args, " "),
		)
	}

	return err.Stderr
}

type repositoryGit struct {
	path             string
	configReadScope  byte
	configWriteScope int
}

/*
NewGitRepository returns a Git Repository for a given path. It does not check that the
path is valid or that the repo is initialialised.

The default config read scope is GitConfigScopeLocal | GitConfigScopeGlobal |
GitConfigScopeSystem. This can be changed by passing a byte as the second
argument, which is useful for local testing.
*/
func NewGitRepository(path string, configReadScope ...byte) Repository {
	var readScope byte = GitConfigScopeLocal | GitConfigScopeGlobal | GitConfigScopeSystem

	if len(configReadScope) > 0 {
		readScope = configReadScope[0]
	}

	return &repositoryGit{
		path: path,

		// Git defaults as of v2.18.0
		configReadScope:  readScope,
		configWriteScope: GitConfigScopeLocal,
	}
}

func (repo *repositoryGit) IsInitialised() bool {
	_, _, _, err := repo.runGitCommandRaw("", "rev-parse")

	return err == nil
}

func (repo *repositoryGit) Init() error {
	_, err := repo.runGitCommand("init")
	return err
}

func (repo *repositoryGit) AllConfig() (map[string]string, error) {
	result, err := repo.runGitCommand("config", repo.configReadScopeArgs(), "--null", "--list")
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
			configMap[kvParts[0]] = ""
		} else {
			configMap[kvParts[0]] = kvParts[1]
		}
	}

	return configMap, nil
}

func (repo *repositoryGit) GetConfig(key string) (string, error) {
	result, err := repo.runGitCommand("config", repo.configReadScopeArgs(), "--get", key)
	if err != nil {
		return "", NewGitConfigErr(err.(*GitCmdErr))
	}

	return result, nil
}

func (repo *repositoryGit) SetConfig(key, value string) error {
	_, err := repo.runGitCommand("config", repo.configWriteScopeArgs(), key, value)
	if err != nil {
		return NewGitConfigErr(err.(*GitCmdErr))
	}

	return nil
}

func (repo *repositoryGit) UnsetConfig(key string) error {
	_, err := repo.runGitCommand("config", repo.configWriteScopeArgs(), "--unset", key)
	if err != nil {
		return NewGitConfigErr(err.(*GitCmdErr))
	}

	return nil
}

func (repo *repositoryGit) GetMetadata(key string) (string, error) {
	id, err := repo.objectID(key)
	if err != nil {
		return "", err
	}

	return repo.runGitCommand("notes", "--ref", gitNotesRef, "show", id)
}

func (repo *repositoryGit) SetMetadata(key, value string) error {
	id, err := repo.objectID(key)
	if err != nil {
		return err
	}

	_, err = repo.runGitCommand("notes", "--ref", gitNotesRef, "add", "-f", id, "-m", value)

	return err
}

func (repo *repositoryGit) runGitCommandRaw(stdin string, args ...string) (string, string, int, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.path
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if stdin != "" {
		cmd.Stdin = bytes.NewBufferString(stdin)
	}

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

func (repo *repositoryGit) runGitCommand(args ...string) (string, error) {
	stdout, stderr, exitCode, err := repo.runGitCommandRaw("", args...)
	if err != nil {
		return "", NewGitCmdErr(stderr, exitCode, args...)
	}

	return stdout, err
}

func (repo *repositoryGit) runGitCommandStdIn(stdin string, args ...string) (string, error) {
	stdout, stderr, exitCode, err := repo.runGitCommandRaw(stdin, args...)
	if err != nil {
		return "", NewGitCmdErr(stderr, exitCode, args...)
	}

	return stdout, err
}

func (repo *repositoryGit) configReadScopeArgs() string {
	args := []string{}

	if (repo.configReadScope & GitConfigScopeLocal) != 0 {
		args = append(args, "--local")
	}
	if (repo.configReadScope & GitConfigScopeGlobal) != 0 {
		args = append(args, "--global")
	}
	if (repo.configReadScope & GitConfigScopeSystem) != 0 {
		args = append(args, "--system")
	}

	return strings.Join(args, " ")
}

func (repo *repositoryGit) configWriteScopeArgs() string {
	switch repo.configWriteScope {
	case GitConfigScopeLocal:
		return "--local"
	case GitConfigScopeGlobal:
		return "--global"
	case GitConfigScopeSystem:
		return "--system"
	}

	return ""
}

func (repo *repositoryGit) objectID(key string) (string, error) {
	return repo.runGitCommandStdIn(key, "hash-object", "--stdin")
}
