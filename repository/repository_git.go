package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/endiangroup/specstack/persistence"
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
	if gitCmdErr.ExitCode == 1 {
		return persistence.ErrNoConfigFound
	}

	return gitCmdErr
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

type Git struct {
	path             string
	configReadScope  int
	configWriteScope int
}

/*
NewGitRepository returns a Git Repository for a given path. It does not
check that the path is valid or that the repo is initialialised.

The default config read scope is GitConfigScopeGlobal . This can be changed
by passing an int in the second argument, which is useful for local
testing.
*/
func NewGitRepository(path string, configReadScope ...int) *Git {
	var readScope int = GitConfigScopeGlobal

	if len(configReadScope) > 0 {
		readScope = configReadScope[0]
	}

	return &Git{
		path:             path,
		configReadScope:  readScope,
		configWriteScope: GitConfigScopeLocal,
	}
}

func (repo *Git) IsInitialised() bool {
	_, _, _, err := repo.runGitCommandRaw(nil, "rev-parse")

	return err == nil
}

func (repo *Git) Init() error {
	_, err := repo.runGitCommand("init")
	return err
}

func (repo *Git) AllConfig() (map[string]string, error) {
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

func (repo *Git) GetConfig(key string) (string, error) {
	result, err := repo.runGitCommand("config", repo.configReadScopeArgs(), "--get", key)
	if err != nil {
		return "", NewGitConfigErr(err.(*GitCmdErr))
	}

	return result, nil
}

func (repo *Git) SetConfig(key, value string) error {
	_, err := repo.runGitCommand("config", repo.configWriteScopeArg(), key, value)
	if err != nil {
		return NewGitConfigErr(err.(*GitCmdErr))
	}

	return nil
}

func (repo *Git) UnsetConfig(key string) error {
	_, err := repo.runGitCommand("config", repo.configWriteScopeArg(), "--unset", key)
	if err != nil {
		return NewGitConfigErr(err.(*GitCmdErr))
	}

	return nil
}

func (repo *Git) GetMetadata(key io.Reader) ([][]byte, error) {
	id, err := repo.objectID(key)
	if err != nil {
		return nil, err
	}

	raw := [][]byte{}

	// Check to see if there's a revision history for this object.
	// If there isn't, we can still check for notes attached directly
	// to the object hash.
	revisions, err := repo.revList(id)
	if err != nil {
		if err := repo.extractJsonMessagesFromObjectId(id, &raw); err != nil {
			return nil, err
		}
	}

	if err := repo.extractJsonMessagesFromRevisionList(revisions, &raw); err != nil {
		return nil, err
	}

	return raw, nil
}

func (repo *Git) extractJsonMessagesFromNote(note string, raw *[][]byte) error {

	for line, value := range strings.Split(note, "\n") {
		decoded := []byte{}
		if err := json.Unmarshal([]byte(value), &decoded); err != nil {
			return fmt.Errorf("failed to parse json from note line %d: %s", line, err)
		}
		*raw = append(*raw, decoded)
	}

	return nil
}

func (repo *Git) extractJsonMessagesFromObjectId(id string, raw *[][]byte) error {
	note, err := repo.runGitCommand("notes", "--ref", gitNotesRef, "show", id)

	// If there is no check and the note recovery fails, it's not an always
	// error; it may mean there are no notes for the id.
	if err != nil {
		if e, ok := err.(*GitCmdErr); ok {
			if strings.HasPrefix(e.Stderr, "error: no note found for object") {
				return nil
			}
		}
		return err
	}

	return repo.extractJsonMessagesFromNote(note, raw)
}

func (repo *Git) extractJsonMessagesFromRevisionList(revisions [][]string, raw *[][]byte) error {
	processedNotes := make(map[string]struct{})
	for _, revision := range revisions {
		if len(revision) > 1 {

			// We want the objects referenced by the commits,
			// in case they're notes. These are in indexes [1:]
			for _, ref := range revision[1:] {

				if _, exists := processedNotes[ref]; exists {
					continue
				}

				if note, err := repo.runGitCommand("notes", "--ref", gitNotesRef, "show", ref); err == nil {
					if err := repo.extractJsonMessagesFromNote(note, raw); err != nil {
						return err
					}
					processedNotes[ref] = struct{}{}
				}
			}
		}
	}
	return nil
}

func (repo *Git) SetMetadata(target io.Reader, value []byte) error {
	id, err := repo.objectID(target)
	if err != nil {
		return err
	}

	encodedValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata value to json: %s", err)
	}

	var note string

	if existingNote, err := repo.runGitCommand("notes", "--ref", gitNotesRef, "show", id); err == nil {
		note = existingNote + "\n" + string(encodedValue)
	} else {
		note = string(encodedValue)
	}

	_, err = repo.runGitCommand("notes", "--ref", gitNotesRef, "add", "-f", id, "-m", note)

	return err
}

func (repo *Git) PrepareMetadataSync() error {
	if err := repo.WriteHookFile("post-commit", "spec git-hook exec post-commit"); err != nil {
		return err
	}

	if err := repo.WriteHookFile("post-update", "spec git-hook exec post-update"); err != nil {
		return err
	}

	return nil
}

func (repo *Git) WriteHookFile(name, command string) error {
	hooksDir, err := repo.gitHooksDirectory()
	if err != nil {
		return err
	}

	path := filepath.Join(hooksDir, name)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}

	content := fmt.Sprintf(
		`#!/bin/sh
# Added by spec command.
%s
`, command)

	return ioutil.WriteFile(path, []byte(content), 0774)
}

func (repo *Git) PullMetadata(from string) error {
	if _, err := repo.runGitCommand("remote", "get-url", from); err != nil {
		return err
	}
	return nil
}

func (repo *Git) PushMetadata(to string) error {
	if _, err := repo.runGitCommand("remote", "get-url", "--push", to); err != nil {
		return err
	}
	return nil
}

/*
revList gets a list of object revisions for a given hash ID.
The information is returned reverse-chronlogically, in N columns:
the first column is the hash of the commit, and the subsequent
(usually blank) columns are the hash of any objects referenced by
the commits.
*/
func (repo *Git) revList(id string) ([][]string, error) {
	output, err := repo.runGitCommand("rev-list", "--all", "--objects", id)
	if err != nil {
		return nil, err
	}

	revisions := [][]string{}

	for _, state := range strings.Split(output, "\n") {
		revision := strings.Split(strings.TrimSpace(state), " ")
		revisions = append(revisions, revision)
	}

	return revisions, nil
}

func (repo *Git) gitHooksDirectory() (string, error) {
	gitDir, err := repo.gitDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(gitDir, "hooks"), nil
}

func (repo *Git) gitDirectory() (string, error) {
	topDir, err := repo.topDirectory()
	if err != nil {
		return "", err
	}
	return filepath.Join(topDir, ".git"), nil
}

func (repo *Git) topDirectory() (string, error) {
	return repo.runGitCommand("rev-parse", "--show-toplevel")
}

func (repo *Git) runGitCommandRaw(stdin io.Reader, args ...string) (string, string, int, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = repo.path
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Stdin = stdin

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
	stdout, stderr, exitCode, err := repo.runGitCommandRaw(nil, args...)
	if err != nil {
		return "", NewGitCmdErr(stderr, exitCode, args...)
	}

	return stdout, err
}

func (repo *Git) runGitCommandStdIn(stdin io.Reader, args ...string) (string, error) {
	stdout, stderr, exitCode, err := repo.runGitCommandRaw(stdin, args...)
	if err != nil {
		return "", NewGitCmdErr(stderr, exitCode, args...)
	}

	return stdout, err
}

func (repo *Git) configReadScopeArgs() string {
	if repo.configReadScope == GitConfigScopeLocal {
		return "--local"
	}
	return ""
}

func (repo *Git) configWriteScopeArg() string {
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

func (repo *Git) objectID(key io.Reader) (string, error) {
	return repo.runGitCommandStdIn(key, "hash-object", "--stdin")
}
