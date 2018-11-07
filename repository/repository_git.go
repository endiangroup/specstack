package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

const (
	notesRef = "refs/notes/specstack"
)

// NewGitConfigErr creates the appropriate typed error for a Git failure, if
// possible.
func NewGitConfigErr(gitCmdErr *GitCmdErr) error {
	switch gitCmdErr.ExitCode {
	case 1:
		return GitConfigMissingSectionKeyErr{gitCmdErr}
	}

	return gitCmdErr
}

// GitConfigMissingSectionKeyErr is returned when a config section is missing
type GitConfigMissingSectionKeyErr struct {
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

// GitCmdErr is an error types for failed Git commands
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
	Path            string
	ConfigNamespace string
}

// NewGit creates a new Git Repository instance from a path and config
// namespace.
func NewGit(path, configNamespace string) Repository {
	return &repositoryGit{
		Path:            path,
		ConfigNamespace: configNamespace,
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

func (repo *repositoryGit) All() (map[string]string, error) {
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

func (repo *repositoryGit) Get(key string) (string, error) {
	return repo.runGitCommand("config", "--get", repo.prefixConfigNamespace(key))
}

func (repo *repositoryGit) Set(key, value string) error {
	_, err := repo.runGitCommand("config", repo.prefixConfigNamespace(key), value)
	return err
}

func (repo *repositoryGit) Unset(key string) error {
	_, err := repo.runGitCommand("config", "--unset", repo.prefixConfigNamespace(key))

	return err
}

func (repo *repositoryGit) GetMetadata(key string) (string, error) {

	id, err := repo.objectID(key)

	if err != nil {
		return "", err
	}

	return repo.runGitCommand("notes", "--ref", notesRef, "show", id)
}

func (repo *repositoryGit) SetMetadata(key, value string) error {

	id, err := repo.objectID(key)

	if err != nil {
		return err
	}

	_, err = repo.runGitCommand("notes", "--ref", notesRef, "add", "-f", id, "-m", value)

	return err
}

func (repo *repositoryGit) runGitCommandRaw(stdin string, args ...string) (string, string, int, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.Path
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

func (repo *repositoryGit) prefixConfigNamespace(key string) string {
	return repo.ConfigNamespace + "." + key
}

func (repo *repositoryGit) trimConfigNamespace(key string) string {
	return strings.TrimPrefix(key, repo.ConfigNamespace+".")
}

func (repo *repositoryGit) objectID(key string) (string, error) {
	return repo.runGitCommandStdIn(key, "hash-object", "--stdin")
}
