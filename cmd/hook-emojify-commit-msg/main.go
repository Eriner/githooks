package main

import (
	"crypto/rand"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"unicode"

	"github.com/eriner/githooks/internal"
)

var fileName string

func init() {
	internal.Init()
	if len(os.Args) != 2 {
		log.Fatalln("error: invalid argument count")
	}
	fileName = os.Args[1]
	if fileName == "" {
		log.Fatalln("error: empty file name provided")
	}
}

func main() {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("error reading commit file: %s\n", err.Error())
	}
	msg := string(dat)

	// validation
	lines := strings.Split(msg, "\n")
	if len(lines) == 0 {
		log.Fatalln("error: commit message is empty")
	}
	if len(lines[0]) < 8 {
		log.Fatalln("error: minimum commit message length is 8 characters")
	}
	if len(lines[0]) > 48 { // 50 - 2, where 50 is max line length and 2 is (emoji) + " "
		log.Fatalf("Summary line is too long! Use less than 48 characters.\n")
	}

	// if there is already an emoji, exit early
	if msg[0] > unicode.MaxASCII {
		log.Println("ok")
		os.Exit(0)
	}

	// pick an emoji by msg contents
	var emoji string
	var emojis map[string]string = map[string]string{
		"init":      `๐งฎ`,
		"fix":       `๐`,
		"bug":       `๐`,
		"ugh":       `๐ฃ`,
		"little":    `๐ค`,
		"hope":      `๐ค`,
		"bump":      `โ`,
		"pray":      `๐`,
		"better":    `๐ช`,
		"think":     `๐ง `,
		"maybe":     `๐คท`,
		"revert":    `๐ฎ`,
		"found":     `๐ต`,
		"refactor":  `๐ช `,
		"christmas": `๐`,
		"santa":     `๐`,
		"update":    `โฌ`,
		"upgrade":   `โฌ`,
		"please":    `๐คฒ`,
		"typo":      `๐ฉน`,
		"trap":      `๐ชค`,
		"cleanup":   `๐งน`,
		"mess":      `๐งน`,
		"debug":     `๐ฉบ`,
		"diagnose":  `๐ฉบ`,
		"garbage":   `๐ฝ`,
		"trash":     `๐ฝ`,
		"rip":       `๐ชฆ`,
		"warning":   `๐ง`,
		"readme":    `๐`,
		"docs":      `๐`,
		"drink":     `๐ฅ`,
		"cheers":    `๐ป`,
		"beer":      `๐บ`,
		"fun":       `๐ก`,
		"ride":      `๐ด`,
		"gas":       `โฝ`,
		"moon":      `๐`,
		"time":      `โณ`,
		"hot":       `๐ก`,
		"yay":       `๐`,
		"hooray":    `๐`,
		"present":   `๐`,
		"ticket":    `๐`,
	}
	for _, str := range strings.Split(lines[0], " ") {
		if emojis[str] != "" {
			emoji = emojis[strings.ToLower(str)]
			break
		}
	}
	if emoji == "" {
		var randomEmojis []string = []string{
			`๐บ`,
			`๐ฅ`,
			`โณ`,
			`๐คฟ`,
			`๐ท`,
			`๐ฑ`,
			`๐ช`,
			`๐ฎ`,
			`๐น`,
			`๐งธ`,
			`โ`,
			`๐ถ`,
			`๐ฆ`,
			`๐`,
			`๐จ`,
			`๐ซ`,
			`๐ช`,
			`๐ง`,
		}
		i, err := rand.Int(rand.Reader, big.NewInt(int64(len(randomEmojis)-1)))
		if err != nil {
			log.Fatalf("error picking random emoji: %s\n", err.Error())
		}
		emoji = randomEmojis[int(i.Int64())]
	}
	msg = emoji + "  " + msg
	if err := ioutil.WriteFile(fileName, []byte(msg), 0o700); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
	log.Println("ok")
}
