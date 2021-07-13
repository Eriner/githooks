package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eriner/githooks/internal"
	"github.com/eriner/githooks/internal/git"
)

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
)

func init() {
	internal.Init()
}

func main() {
	var bigFiles []string
	nf, err := git.NewStagedFiles()
	if err != nil {
		log.Fatalln(err)
	}
	for _, fileName := range nf {
		s, err := os.Stat(fileName)
		if err != nil {
			log.Fatalf("error: unable to get file size")
		}
		if s.Size() > (20 * MiB) {
			bigFiles = append(bigFiles, fileName)
		}
	}
	if len(bigFiles) == 0 {
		log.Println("no big files, all good!")
		os.Exit(0)
	}
	log.Println("The following files seem rather large, are you sure that you meant to add them?")
	bigFilesStr := "  * " + strings.Join(bigFiles, "\n  * ")
	log.Println(bigFilesStr)
	log.Println("If you don't want to add these files, abort now with Ctrl+C")
	log.Println("Otherwise, confirm by typing \"ok\"")
	fmt.Print("> ")
	out, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if len(out) < 1 {
		log.Println("aborting.")
		os.Exit(1)
	}
	if strings.TrimSpace(strings.ToLower(out)) != "ok" {
		log.Println("aborting.")
		os.Exit(1)
	}
	os.Exit(0)
}
