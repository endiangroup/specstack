package repository

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

const (
	ScopeLocal  = 1
	ScopeSystem = 2
	ScopeGlobal = 4
)

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

type Git struct {
	Path             string
	ConfigReadScope  byte
	ConfigWriteScope int
}

func NewGit(path string) *Git {
	return &Git{
		Path: path,

		// Git defaults as of v2.18.0
		ConfigReadScope:  ScopeLocal | ScopeGlobal | ScopeSystem,
		ConfigWriteScope: ScopeLocal,
	}
}

func (repo *Git) SetConfigReadScope(scope byte) {
	repo.ConfigReadScope = scope
}
func (repo *Git) SetConfigWriteScope(scope int) {
	repo.ConfigWriteScope = scope
}

func (repo *Git) IsInitialised() bool {
	_, _, _, err := repo.runGitCommandRaw("rev-parse")

	return err == nil
}

func (repo *Git) Init() error {
	_, err := repo.runGitCommand("init")
	return err
}

func (repo *Git) All() (map[string]string, error) {
	result, err := repo.runGitCommand("config", repo.configReadScope(), "--null", "--list")
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

func (repo *Git) Get(key string) (string, error) {
	result, err := repo.runGitCommand("config", repo.configReadScope(), "--get", key)
	if err != nil {
		return "", NewGitConfigErr(err.(*GitCmdErr))
	}

	return result, nil
}

func (repo *Git) Set(key, value string) error {
	_, err := repo.runGitCommand("config", repo.configWriteScope(), key, value)
	if err != nil {
		return NewGitConfigErr(err.(*GitCmdErr))
	}

	return nil
}

func (repo *Git) Unset(key string) error {
	_, err := repo.runGitCommand("config", repo.configWriteScope(), "--unset", key)
	if err != nil {
		return NewGitConfigErr(err.(*GitCmdErr))
	}

	return nil
}

func (repo *Git) runGitCommandRaw(args ...string) (string, string, int, error) {
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

func (repo *Git) runGitCommand(args ...string) (string, error) {
	stdout, stderr, exitCode, err := repo.runGitCommandRaw(args...)
	if err != nil {
		return "", NewGitCmdErr(stderr, exitCode, args...)
	}

	return stdout, err
}

func (repo *Git) configReadScope() string {
	args := []string{}

	if (repo.ConfigReadScope & ScopeLocal) != 0 {
		args = append(args, "--local")
	}
	if (repo.ConfigReadScope & ScopeGlobal) != 0 {
		args = append(args, "--global")
	}
	if (repo.ConfigReadScope & ScopeSystem) != 0 {
		args = append(args, "--system")
	}

	return strings.Join(args, " ")
}

func (repo *Git) configWriteScope() string {
	switch repo.ConfigWriteScope {
	case ScopeLocal:
		return "--local"
	case ScopeGlobal:
		return "--global"
	case ScopeSystem:
		return "--system"
	}

	return ""
}
