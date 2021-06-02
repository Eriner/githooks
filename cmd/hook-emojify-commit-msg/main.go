package main

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

var fileName string

func init() {
	log.Default().SetPrefix(fmt.Sprintf("%s: ", filepath.Base(os.Args[0])))
	log.Default().SetFlags(0)
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
	log.Println("commit looks good!")
	if err := ioutil.WriteFile(fileName, []byte(msg), 0o700); err != nil {
		log.Fatalf("error: %s\n", err.Error())
	}
}
