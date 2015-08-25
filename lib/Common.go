package lib

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	yes_text = "Yeah"
	no_text  = "Nope"
)

type (
	Part struct {
		Func    func(incoming telebot.Message)
		ToBeRun bool
	}

	Parts []*Part

	option struct {
		Text             string   `json:"text"`
		Feedback         []string `json:"feedback"`
		FollowUpQuestion []string `json:"followUp"`
	}

	Details struct {
		Bot  *telebot.Bot
		User telebot.User
	}

	Process interface {
		GetParts() Parts
	}
)

func loadLanguage(lang interface{}) {
	file, _ := os.Open("./languageConfig.json")
	decoder := json.NewDecoder(file)

	log.Println("Loaded language config file", file.Name())

	err := decoder.Decode(&lang)
	if err != nil {
		log.Panic("could not load language ", err.Error())
	}
}

func getLangText(opts []string) string {
	return opts[rand.Intn(len(opts))]
}

func makeKeyBoard(keys ...string) [][]string {
	keyboard := [][]string{}
	for i := range keys {
		keyboard = append(keyboard, []string{keys[i]})
	}
	return keyboard
}

func (d *Details) send(msg string) {
	d.Bot.SendMessage(
		d.User,
		msg,
		nil,
	)
	return
}

func (d *Details) sendWithKeyboard(msg string, keyboard [][]string) {
	d.Bot.SendMessage(
		d.User,
		msg,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply:      true,
				CustomKeyboard:  keyboard,
				ResizeKeyboard:  false,
				OneTimeKeyboard: true,
			},
		},
	)
	return
}

func Run(m telebot.Message, p Parts) {
	for i := range p {
		log.Println("checking ", i, "of", len(p), "run?", p[i].ToBeRun)
		if p[i].ToBeRun {
			log.Println("running ", i)
			p[i].ToBeRun = false
			p[i].Func(m)
			log.Println("has run ", i)
			return
		}
	}
}