package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/eriner/githooks/internal/git"
	"github.com/gabriel-vasile/mimetype"
	"github.com/montanaflynn/stats"
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
	stagedFiles, err := git.GetStagedFiles()
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
		/*
			if suspect[f] {
				continue
			}
		*/

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

func isEntropic(file string) (bool, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return false, err
	}
	mime := mimetype.Detect(dat)

	if mime.Is("text/plain") {
		// Split the lines and calculate the Shanon entropy for each line.
		// If there are extreme outliers, it is possible that there are
		// hardcoded secrets.
		lines := strings.Split(string(dat), "\n")
		var scores []float64
		for _, line := range lines {
			scores = append(scores, float64(shannon(line)))
		}
		scores = append(scores, float64(shannon(string(dat))))

		mean, err := stats.Mean(scores)
		if err != nil {
			return false, fmt.Errorf("error calculating shannon score mean: %w", err)
		}
		// log.Printf("file %q mean %v", file, mean)
		if mean > 500 {
			// If the mean is more than 500, it is likely the file is a certificate or
			// some other truly random data.
			return true, nil
		}

		// NOTE: my idea to detect one-off hardcoded secrets using statistical outliers still
		// needs some work here. Given a large enough file, there are bound to be outliers in
		// the fourth quartile. But this alone does not signal entropy. Will come back to this
		// at some point. The mean check should be enough to catch private keys.
		/*

			// If we don't have enough datapoints, calculating statistical outliers is pointless
			// and we can return early.
			if len(scores) < 24 {
				return false, nil
			}

			// For now, we're considering any file with extreme outliers above the mean to be entropic.
			// This will have more false-positives, but this also ensures that we don't miss anything.
			outliers, err := stats.QuartileOutliers(scores)
			if err != nil {
				return false, fmt.Errorf("error calculating shannon score outliers: %w", err)
			}
			for _, outlier := range outliers.Extreme {
				if outlier > mean {
					log.Println("exteme outlier: " + file)
					return true, nil
				}
			}
		*/
	}
	return false, nil
}

func shannon(s string) int {
	// for the sake of not triggering on every source-langage file...
	charFreq := make(map[rune]float64)
	for _, i := range s {
		charFreq[i]++
	}
	var t float64
	for _, freq := range charFreq {
		f := freq / float64(len(s))
		t += f * math.Log2(f)
	}
	return int(math.Ceil(t*-1)) * len(s)
}
