// package git is not meant to replace github.com/go-git/go-git,
// rather I just don't want the extra dependency until I need it.
// All functions imply that ${PWD}/.git exists, and that the `git`
// binary is available.
package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func GetStagedFiles() ([]string, error) {
	return execGitCmd([]string{"git", "diff", "--staged", "--name-only", "--no-color"})
}

// GetNewSGetNewStagedFiles is similar to GetStagedFiles, but only returns the new
// files staged for commit, not modified files.
func GetNewStagedFiles() ([]string, error) {
	return execGitCmd([]string{"git", "diff", "--staged", "--name-only", "--no-color", "--diff-filter=A", "HEAD"})
}

// execGitCmd returns a slice of files as sourced from running `git` commands.
func execGitCmd(cmd []string) ([]string, error) {
	if len(cmd) < 2 {
		return nil, errors.New("invalid command passed to execGitCmd")
	}
	if cmd[0] != "git" {
		return nil, errors.New("non-git command passed to execGitCmd")
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	c.Stdout = stdout
	c.Stderr = stderr
	err := c.Run()
	if err != nil {
		return nil, fmt.Errorf("error running git diff: %w", err)
	}
	rawLines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	var out []string
	for _, str := range rawLines {
		if strings.TrimSpace(str) == "" {
			continue
		}
		out = append(out, str)
	}
	return out, nil
}
