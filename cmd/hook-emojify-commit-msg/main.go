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
		"init":      `🧮`,
		"fix":       `🖍`,
		"bug":       `🐛`,
		"ugh":       `😣`,
		"little":    `🤏`,
		"hope":      `🤞`,
		"bump":      `✊`,
		"pray":      `🙏`,
		"better":    `💪`,
		"think":     `🧠`,
		"maybe":     `🤷`,
		"revert":    `👮`,
		"found":     `🕵`,
		"refactor":  `🪠`,
		"christmas": `🎄`,
		"santa":     `🎅`,
		"update":    `⬆`,
		"upgrade":   `⬆`,
		"please":    `🤲`,
		"typo":      `🩹`,
		"trap":      `🪤`,
		"cleanup":   `🧹`,
		"mess":      `🧹`,
		"debug":     `🩺`,
		"diagnose":  `🩺`,
		"garbage":   `🚽`,
		"trash":     `🚽`,
		"rip":       `🪦`,
		"warning":   `🚧`,
		"readme":    `📚`,
		"docs":      `📚`,
		"drink":     `🥃`,
		"cheers":    `🍻`,
		"beer":      `🍺`,
		"fun":       `🎡`,
		"ride":      `🛴`,
		"gas":       `⛽`,
		"moon":      `🚀`,
		"time":      `⏳`,
		"hot":       `🌡`,
		"yay":       `🎉`,
		"hooray":    `🎊`,
		"present":   `🎁`,
		"ticket":    `🎟`,
	}
	for _, str := range strings.Split(lines[0], " ") {
		if emojis[str] != "" {
			emoji = emojis[strings.ToLower(str)]
			break
		}
	}
	if emoji == "" {
		var randomEmojis []string = []string{
			`😺`,
			`🥋`,
			`⛳`,
			`🤿`,
			`🛷`,
			`🎱`,
			`🪄`,
			`🎮`,
			`🕹`,
			`🧸`,
			`♟`,
			`🕶`,
			`📦`,
			`📝`,
			`🔨`,
			`🔫`,
			`🪃`,
			`🔧`,
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
