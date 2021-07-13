package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

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
	var jsonCt int
	for _, f := range stagedFiles {
		if strings.HasSuffix(f, ".json") {
			jsonCt++
			dat, err := ioutil.ReadFile(f)
			if err != nil {
				log.Fatalf("error reading file %q: %s\n", f, err.Error())
			}
			if !json.Valid(dat) {
				log.Fatalln("json file did not parse as valid json")
			}
		}
	}
	if jsonCt > 0 {
		log.Println("ok")
		os.Exit(0)
	}
}
