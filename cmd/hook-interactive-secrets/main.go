package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/eriner/githooks/internal/git"
)

var (
	suspectFileExtensions []string = []string{".key", ".pem"}
	whitelistedExtensions []string = []string{".sig", ".pub", ".sum"}
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
	if len(stagedFiles) == 0 {
		log.Fatalln("no files staged for commit")
	}

	// suspect files are the ones we intend to prompt the user to be certain
	suspect := make(map[string]bool)
	for _, f := range stagedFiles {
		suspect[f] = false
		var ee bool
		for _, suffix := range whitelistedExtensions {
			if strings.HasSuffix(f, suffix) {
				ee = true
			}
		}
		if ee {
			continue
		}

		// Check for suspect extensions
		for _, suffix := range suspectFileExtensions {
			if strings.HasSuffix(f, suffix) {
				suspect[f] = true
			}
		}
		if suspect[f] {
			continue
		}

		// Check Shannon entropy
		entropic, err := isEntropic(f)
		if err != nil {
			log.Fatalf("error checking file entropy: %s", err.Error())
		}
		if entropic {
			suspect[f] = true
			continue
		}

	}

	// quick check to see if we can exit early
	var suspectCt int
	for _, isSuspect := range suspect {
		if isSuspect {
			suspectCt++
		}
	}
	if suspectCt == 0 {
		log.Println("no secrets staged for commit, all good!")
		os.Exit(0)
	}

	log.Println()
	log.Println("-----------------------------------------------------")
	log.Println("|     !!!     Secrets Exposure Warning      !!!     |")
	log.Println("-----------------------------------------------------")
	log.Println()
	log.Println("The following files are suspect, and may contain secrets:")
	for file, check := range suspect {
		if !check {
			continue
		}
		log.Println("  - " + file)
	}
	log.Println()
	log.Println("If these files are sensitive or contain secrets, abort now with Ctrl+C")
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
}

const (
	b64Charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/="
	hexCharset string = "0123456789abcdefABCEDF"
)

// limitCharset returns s, limited to characters in charset, with a limiter cap
// of 20 characters. After there are 20 chararcters that match the charset, the rest
// are added *unconditionally*.
// Not to be used as a filter or for validation routines.
func limitCharset(s string, charset string) string {
	var out string
	for _, ch := range s {
		if len(out) > 20 {
			out += string(ch)
		}
		if x := strings.IndexRune(charset, ch); x != -1 {
			out += string(ch)
		}
	}
	return out
}

func isEntropic(file string) (bool, error) {
	dat, err := git.DiffOfFile(file)
	if err != nil {
		return false, err
	}
	lines := strings.Split(dat, "\n")
	for _, line := range lines {
		for _, word := range strings.Split(line, " ") {
			b64Shannon := float64(shannon(limitCharset(word, b64Charset)))
			if b64Shannon > 4.5 {
				log.Printf("entropic string in %q found: %q\n", file, word)
				return true, nil
			}
			hexShannon := float64(shannon(limitCharset(word, hexCharset)))
			if hexShannon > 3.5 {
				log.Printf("entropic string in %q found: %q\n", file, word)
				return true, nil
			}
		}
	}
	return false, nil
}

func shannon(s string) float64 {
	if s == "" {
		return 0
	}
	charFreq := make(map[rune]float64)
	for _, i := range s {
		charFreq[i]++
	}
	var t float64
	for _, freq := range charFreq {
		f := freq / float64(len(s))
		if f > 0 {
			t += -1 * (f * math.Log2(f))
		}
	}
	return t
}
