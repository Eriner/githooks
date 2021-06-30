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
	"path/filepath"
	"strings"
)

// DiffOfFile returns `git diff` of files, minus the diff header.
func DiffOfFile(file string) (string, error) {
	if strings.Index(file, " ") != -1 {
		return "", errors.New("invalid file provided to GetDiffOfFile")
	}
	f, err := filepath.Abs(filepath.Clean(file))
	if err != nil {
		return "", fmt.Errorf("error getting file absolute path: %q", f)
	}
	s, err := execGitCmd([]string{"git", "diff", "--staged", f})
	if err != nil {
		return "", err
	}
	// We have to lop off the 4-line header, as it contains the git short-hashes which
	// are not a part of the commit and guaranteed to have high-entropy. Maybe there
	// is a git flag for this. /shrug
	ss := strings.Split(s, "\n")
	if len(ss) < 4 {
		return "", errors.New("git diff missing header: 4-line header expected in output, but was missing")
	}
	return strings.Join(ss[3:], "\n"), nil
}

func StagedFiles() ([]string, error) {
	s, err := execGitCmd([]string{"git", "diff", "--staged", "--name-only", "--no-color"})
	return sliceFromStringList(s), err
}

// NewStagedFiles is similar to StagedFiles, but only returns the new
// files staged for commit, not modified files.
func NewStagedFiles() ([]string, error) {
	s, err := execGitCmd([]string{"git", "diff", "--staged", "--name-only", "--no-color", "--diff-filter=A", "HEAD"})
	return sliceFromStringList(s), err
}

// execGitCmd returns a slice of files as sourced from running `git` commands.
func execGitCmd(cmd []string) (string, error) {
	if len(cmd) < 2 {
		return "", errors.New("invalid command passed to execGitCmd")
	}
	if cmd[0] != "git" {
		return "", errors.New("non-git command passed to execGitCmd")
	}
	c := exec.Command(cmd[0], cmd[1:]...)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	c.Stdout = stdout
	c.Stderr = stderr
	err := c.Run()
	if err != nil {
		return "", fmt.Errorf("error running git diff: %w", err)
	}
	return stdout.String(), nil
}

func sliceFromStringList(s string) []string {
	rawLines := strings.Split(strings.TrimSpace(s), "\n")
	var out []string
	for _, str := range rawLines {
		if strings.TrimSpace(str) == "" {
			continue
		}
		out = append(out, str)
	}
	return out
}
