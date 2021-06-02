package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/eriner/githooks/internal/git"
)

func init() {
	log.Default().SetPrefix(fmt.Sprintf("%s: ", filepath.Base(os.Args[0])))
	log.Default().SetFlags(0)
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
		log.Println("all json files parsed as valid json, all good!")
	}
}
