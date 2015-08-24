package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	//commands
	yostart_command = "/start"
	yobg_command    = "/yobg"
	yolow_command   = "/yolow"
	yomove_command  = "/yomove"
	yofood_command  = "/yofood"

	typing_action = "typing"

	note_json = `{"type":"note","user":"%s","data":{"text":"%s"},"time":"%s"}`
)

//our bot
var fBot *fantailBot

type (
	lang struct {
		Greetings []string `json:"greet"`
		Goodbyes  []string `json:"goodbye"`
		Yes       []string `json:"yes"`
		No        []string `json:"no"`
		Bg        struct {
			Comment  string `json:"comment"`
			Question string `json:"question"`
			Above    option `json:"above"`
			In       option `json:"in"`
			Below    option `json:"below"`
			Thank    string `json:"thank"`
		} `json:"bg"`
		Move struct {
			Comment  string `json:"comment"`
			Question string `json:"question"`
			Ex1      string `json:"typeOne"`
			Ex2      string `json:"typeTwo"`
			Ex3      string `json:"typeThree"`
			Ex4      string `json:"typeFour"`
			Thank    string `json:"thank"`
		} `json:"move"`
		Food struct {
			Comment  string `json:"comment"`
			Question string `json:"question"`
			Snack    string `json:"snack"`
			Meal     string `json:"meal"`
			Other    string `json:"other"`
			Thank    string `json:"thank"`
		} `json:"food"`
		Low struct {
			Comment  string `json:"comment"`
			Question string `json:"question"`
			Good     option `json:"good"`
			NotGood  option `json:"notGood"`
			Other    option `json:"other"`
			Thank    string `json:"thank"`
		} `json:"low"`
	}

	option struct {
		Text     string   `json:"text"`
		Feedback []string `json:"feedback"`
		FollowUp []string `json:"followUp"`
	}

	fantailBot struct {
		process
		bot *telebot.Bot
		*lang
	}
	process interface {
		getBot() *telebot.Bot
		getLanguage() *lang
		isRunning(string) bool
		runningName() string
		addPart(*part)
		run(telebot.Message)
	}
	part struct {
		fn    func(msg telebot.Message)
		toRun bool
	}
)

func loadLanguage() *lang {
	file, _ := os.Open("./languageConfig.json")
	decoder := json.NewDecoder(file)
	var language lang
	err := decoder.Decode(&language)
	if err != nil {
		log.Panic("could not load language ", err.Error())
	}
	return &language
}

func getfBot() *fantailBot {
	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("$BOT_TOKEN must be set")
	}

	bot, err := telebot.NewBot(botToken)
	if err != nil {
		return nil
	}
	return &fantailBot{bot: bot, lang: loadLanguage()}
}

func (this *fantailBot) isCurrentlyRunning(processName string) bool {
	if this.process != nil && this.process.runningName() == processName {
		return true
	}
	return false
}

func main() {

	fBot = getfBot()

	messages := make(chan telebot.Message)

	go fBot.bot.Listen(messages, 1*time.Second)

	for msg := range messages {

		log.Println("incomming msg", msg.Text)
		if fBot.process != nil {
			log.Println("running process ", fBot.process.runningName())
		}
		if strings.Contains(msg.Text, yostart_command) {
			//show all options
			b := newBasics(fBot, yostart_command)
			b.options(msg)
		} else if strings.Contains(msg.Text, yobg_command) || fBot.isCurrentlyRunning(yobg_command) {
			if strings.Contains(msg.Text, yobg_command) {
				b := newBasics(fBot, yobg_command)
				b.addPart(&part{fn: b.bg, toRun: true})
				b.addPart(&part{fn: b.bgFeedback, toRun: true})
				b.addPart(&part{fn: b.seeYou, toRun: true})
				fBot.setProcess(b)
			}

			fBot.process.run(msg)
		} else if strings.Contains(msg.Text, yomove_command) || fBot.isCurrentlyRunning(yomove_command) {
			if strings.Contains(msg.Text, yomove_command) {
				b := newBasics(fBot, yomove_command)
				b.addPart(&part{fn: b.yoMove, toRun: true})
				b.addPart(&part{fn: b.seeYou, toRun: true})
				fBot.setProcess(b)
			}

			fBot.process.run(msg)
		} else if strings.Contains(msg.Text, yofood_command) || fBot.isCurrentlyRunning(yofood_command) {
			if strings.Contains(msg.Text, yofood_command) {
				b := newBasics(fBot, yofood_command)
				b.addPart(&part{fn: b.yoFood, toRun: true})
				b.addPart(&part{fn: b.seeYou, toRun: true})
				fBot.setProcess(b)
			}

			fBot.process.run(msg)
		} else if strings.Contains(msg.Text, yolow_command) || fBot.isCurrentlyRunning(yolow_command) {
			if strings.Contains(msg.Text, yolow_command) {
				b := newBasics(fBot, yolow_command)
				b.addPart(&part{fn: b.low, toRun: true})
				b.addPart(&part{fn: b.lowFeedBack, toRun: true})
				b.addPart(&part{fn: b.seeYou, toRun: true})
				fBot.setProcess(b)
			}
			fBot.process.run(msg)
		}

	}
}

type basics struct {
	fBot  *fantailBot
	parts []*part
	name  string
}

func newBasics(fBot *fantailBot, name string) *basics {
	return &basics{fBot: fBot, name: name}
}

func (this *fantailBot) setProcess(p *basics) {
	this.process = nil
	this.process = p
	return
}

func (this *basics) seeYou(msg telebot.Message) {
	this.getBot().SendMessage(
		msg.Chat,
		this.getLanguage().Goodbyes[rand.Intn(len(this.getLanguage().Goodbyes))],
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: true,
				CustomKeyboard: [][]string{
					[]string{this.getLanguage().Goodbyes[rand.Intn(len(this.getLanguage().Goodbyes))]},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		})
	return
}

func (this *basics) pause(msg telebot.Message) {
	this.getBot().SendChatAction(msg.Chat, typing_action)
	time.Sleep(2 * time.Second)
	return
}

func (this *basics) addPart(p *part) {
	this.parts = append(this.parts, p)
	return
}

func (this *basics) isRunning(name string) bool {
	return this.name == name
}

func (this *basics) runningName() string {
	return this.name
}

func (this *basics) run(m telebot.Message) {
	for i := range this.parts {
		log.Println("checking ", i, "of", len(this.parts), "run?", this.parts[i].toRun)
		if this.parts[i].toRun {
			log.Println("running ", i)
			this.parts[i].toRun = false
			this.parts[i].fn(m)
			log.Println("has run ", i)
			return
		}
	}
}

func (this *basics) getBot() *telebot.Bot {
	return this.fBot.bot
}

func (this *basics) getLanguage() *lang {
	return this.fBot.lang
}

func (this *basics) options(msg telebot.Message) {
	this.getBot().SendMessage(msg.Chat,
		fmt.Sprintf("%s %s! What can we do for you?", this.getLanguage().Greetings[rand.Intn(len(this.getLanguage().Greetings))], msg.Chat.FirstName),
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: true,
				CustomKeyboard: [][]string{
					[]string{yobg_command},
					[]string{yofood_command},
					[]string{yomove_command},
					[]string{yolow_command},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: false,
			},
		})
	return
}

func (this *basics) yesNoOpts() *telebot.SendOptions {
	return &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			ForceReply: true,
			CustomKeyboard: [][]string{
				[]string{this.getLanguage().Yes[rand.Intn(len(this.getLanguage().Yes))]},
				[]string{this.getLanguage().No[rand.Intn(len(this.getLanguage().No))]},
			},
			ResizeKeyboard:  true,
			OneTimeKeyboard: true,
		},
	}
}

func (this *basics) bg(msg telebot.Message) {
	this.getBot().SendMessage(msg.Chat,
		this.getLanguage().Bg.Question,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: true,
				CustomKeyboard: [][]string{
					[]string{this.getLanguage().Bg.Above.Text},
					[]string{this.getLanguage().Bg.In.Text},
					[]string{this.getLanguage().Bg.Below.Text},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		})
	return
}

func (this *basics) bgFeedback(msg telebot.Message) {
	switch {
	case msg.Text == this.getLanguage().Bg.Above.Text:
		this.getBot().SendMessage(msg.Chat, this.getLanguage().Bg.Above.Feedback[rand.Intn(len(this.getLanguage().Bg.Above.Feedback))], nil)
		this.pause(msg)
		this.getBot().SendMessage(msg.Chat, this.getLanguage().Bg.Above.FollowUp[rand.Intn(len(this.getLanguage().Bg.Above.FollowUp))], this.yesNoOpts())
		return
	case msg.Text == this.getLanguage().Bg.In.Text:
		this.getBot().SendMessage(msg.Chat, this.getLanguage().Bg.In.Feedback[rand.Intn(len(this.getLanguage().Bg.In.Feedback))], nil)
		this.pause(msg)
		this.getBot().SendMessage(msg.Chat, this.getLanguage().Bg.In.FollowUp[rand.Intn(len(this.getLanguage().Bg.In.FollowUp))], this.yesNoOpts())
		return
	case msg.Text == this.getLanguage().Bg.Below.Text:
		this.getBot().SendMessage(msg.Chat, this.getLanguage().Bg.Below.Feedback[rand.Intn(len(this.getLanguage().Bg.Below.Feedback))], nil)
		this.pause(msg)
		this.getBot().SendMessage(msg.Chat, this.getLanguage().Bg.Below.FollowUp[rand.Intn(len(this.getLanguage().Bg.Below.FollowUp))], this.yesNoOpts())
		return
	}
	return
}

func (this *basics) low(msg telebot.Message) {
	this.getBot().SendMessage(msg.Chat,
		this.getLanguage().Low.Question,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: true,
				CustomKeyboard: [][]string{
					[]string{this.getLanguage().Low.Good.Text},
					[]string{this.getLanguage().Low.NotGood.Text},
					[]string{this.getLanguage().Low.Other.Text},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		})
	return
}

func (this *basics) lowFeedBack(msg telebot.Message) {
	switch {
	case msg.Text == this.getLanguage().Low.Good.Text:
		this.getBot().SendMessage(
			msg.Chat,
			this.getLanguage().Low.Good.Feedback[rand.Intn(len(this.getLanguage().Low.Good.Feedback))],
			nil,
		)
		this.pause(msg)
		this.getBot().SendMessage(
			msg.Chat,
			this.getLanguage().Low.Good.Feedback[rand.Intn(len(this.getLanguage().Low.Good.Feedback))],
			this.yesNoOpts(),
		)
		return
	case msg.Text == this.getLanguage().Low.NotGood.Text:
		this.getBot().SendMessage(
			msg.Chat,
			this.getLanguage().Low.NotGood.Feedback[rand.Intn(len(this.getLanguage().Low.NotGood.Feedback))],
			nil,
		)
		this.pause(msg)
		this.getBot().SendMessage(
			msg.Chat,
			this.getLanguage().Low.NotGood.FollowUp[rand.Intn(len(this.getLanguage().Low.NotGood.FollowUp))],
			this.yesNoOpts(),
		)
		return
	case msg.Text == this.getLanguage().Low.Other.Text:
		this.getBot().SendMessage(
			msg.Chat,
			this.getLanguage().Low.Other.Feedback[rand.Intn(len(this.getLanguage().Low.Other.Feedback))],
			nil,
		)
		this.pause(msg)
		this.getBot().SendMessage(
			msg.Chat,
			this.getLanguage().Low.Other.FollowUp[rand.Intn(len(this.getLanguage().Low.Other.FollowUp))],
			this.yesNoOpts(),
		)
		return
	}
	return
}

func (this *basics) yoFood(msg telebot.Message) {
	this.getBot().SendMessage(msg.Chat,
		this.getLanguage().Food.Question,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: true,
				CustomKeyboard: [][]string{
					[]string{this.getLanguage().Food.Snack},
					[]string{this.getLanguage().Food.Meal},
					[]string{this.getLanguage().Food.Other},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		})
	return
}

func (this *basics) yoMove(msg telebot.Message) {
	this.getBot().SendMessage(msg.Chat,
		this.getLanguage().Move.Question,
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: true,
				CustomKeyboard: [][]string{
					[]string{this.getLanguage().Move.Ex1},
					[]string{this.getLanguage().Move.Ex2},
					[]string{this.getLanguage().Move.Ex3},
					[]string{this.getLanguage().Move.Ex4},
				},
				ResizeKeyboard:  true,
				OneTimeKeyboard: true,
			},
		})
	return
}
