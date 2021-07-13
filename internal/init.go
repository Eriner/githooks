package internal

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Init is a shared initilizer for all of my githooks
func Init() {
	// Each tool must print its name
	log.Default().SetPrefix(fmt.Sprintf("%s: ", filepath.Base(os.Args[0])))
	log.Default().SetFlags(0)
}
