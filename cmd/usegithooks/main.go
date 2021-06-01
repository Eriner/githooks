package main

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

const (
	templatesDir string = ".git-templates"
	dest         string = ".git/hooks"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error getting user home dir: %s\n", err.Error())
	}
	hooksDir := filepath.Join(home, templatesDir, "hooks")
	_, err = ioutil.ReadDir(dest)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("error reading \".git/hooks\" directory: %s\n", err.Error())
	}
	_ = os.Remove(dest + "-old")
	if err := os.Rename(dest, dest+"-old"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("error moving existing hooks directory: %s\n", err.Error())
	}
	if err := os.Symlink(hooksDir, dest); err != nil {
		log.Fatalf("error symlinking %q to %q: %s\n", hooksDir, dest, err.Error())
	}
}
