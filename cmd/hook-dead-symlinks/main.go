package main

import (
	"errors"
	"log"
	"os"

	"github.com/eriner/githooks/internal"
	"github.com/eriner/githooks/internal/git"
)

func init() {
	internal.Init()
}

func main() {
	stagedFiles, err := git.StagedFiles()
	if err != nil {
		log.Fatalf("error getting staged files: %s\n", err.Error())
	}
	var symlinkCt int
	for _, f := range stagedFiles {
		s, err := os.Lstat(f)
		if err != nil {
			if _, ok := err.(*os.PathError); ok {
				// can occur when a staged file is deleted
				continue
			}
			log.Fatalf("error stating file %q: %s\n", f, err.Error())
		}
		if s.Mode()&os.ModeSymlink != 0 {
			symlinkCt++
			pointingToFile, err := os.Readlink(s.Name())
			if err != nil {
				log.Fatalf("symlink %q points to dead link at %q\n", f, pointingToFile)
			}
			if _, err = os.Stat(pointingToFile); errors.Is(err, os.ErrNotExist) {
				log.Fatalf("symlink %q points to dead link at %q\n", f, pointingToFile)
			}
		}
	}
	if symlinkCt > 0 {
		log.Println("no dead symlinks, all clear!")
	}
}
