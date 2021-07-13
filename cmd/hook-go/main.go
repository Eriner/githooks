package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/eriner/githooks/internal"
)

func init() {
	internal.Init()
}

func main() {
	if _, err := os.Stat("go.mod"); errors.Is(err, os.ErrNotExist) {
		// not a Go project
		os.Exit(0)
	}
	commands := []string{
		"go vet ./...",
		"go fmt ./...",
		"go test ./...",
	}
	for _, c := range commands {
		if err := cmd(c); err != nil {
			log.Fatalf("error running %q: %s\n", c, err.Error())
		}
	}
	log.Println("go vet, fmt, and test look good!")
}

func cmd(command string) error {
	if command == "" {
		return errors.New("exec() requires a command")
	}
	s := strings.Split(command, " ")
	var cmd *exec.Cmd
	if len(s) == 1 {
		cmd = exec.Command(s[0])
	} else {
		cmd = exec.Command(s[0], s[1:]...)
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	if cmd.ProcessState.ExitCode() != 0 {
		return fmt.Errorf("exited with status code %d", cmd.ProcessState.ExitCode())
	}
	return nil
}
