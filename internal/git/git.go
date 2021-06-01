// package git is not meant to replace github.com/go-git/go-git,
// rather I just don't want the extra dependency until I need it.
// All functions imply that ${PWD}/.git exists, and that the `git`
// binary is available.
package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func GetStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
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
