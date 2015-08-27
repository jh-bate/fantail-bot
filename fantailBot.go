package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
	"github.com/jh-bate/fantail-bot/lib"
)

const (
	//commands
	yostart_command = "/hey"
	yobg_command    = "/bg"
	yolow_command   = "/low"
)

type (
	fantailBot struct {
		bot     *telebot.Bot
		process lib.Process
	}
)

func newFantailBot() *fantailBot {
	botToken := os.Getenv("BOT_TOKEN")

	if botToken == "" {
		log.Fatal("$BOT_TOKEN must be set")
	}

	bot, err := telebot.NewBot(botToken)
	if err != nil {
		return nil
	}
	return &fantailBot{bot: bot}
}

func (this *fantailBot) showOptions(usr telebot.User) {
	this.bot.SendMessage(
		usr,
		fmt.Sprintf("Hey %s", usr.FirstName),
		nil)
	this.bot.SendMessage(
		usr,
		"Sorry but at this time you can only choose from the options below ...",
		nil)
	this.bot.SendMessage(
		usr,
		fmt.Sprintf("Select one and we can get this party started!"),
		&telebot.SendOptions{
			ReplyMarkup: telebot.ReplyMarkup{
				ForceReply: false,
				CustomKeyboard: [][]string{
					[]string{yobg_command},
					[]string{yolow_command},
				},
				ResizeKeyboard:  false,
				OneTimeKeyboard: false,
			},
		})
	return
}

func (this *fantailBot) setRunning(p lib.Process) {
	this.process = nil
	this.process = p
	return
}

func (this *fantailBot) isRunning() bool {
	return this.process != nil && this.process.CanRun()
}

func (this *fantailBot) startLow(usr telebot.User) {
	log.Println("init LOW setup")
	this.setRunning(lib.NewLow(&lib.Details{Bot: this.bot, User: usr}))
}

func (this *fantailBot) startQuickBg(usr telebot.User) {
	log.Println("init QBG setup")
	this.setRunning(lib.NewQuickBg(&lib.Details{Bot: this.bot, User: usr}))
}

func main() {

	fBot := newFantailBot()
	messages := make(chan telebot.Message)
	fBot.bot.Listen(messages, 1*time.Second)

	for msg := range messages {

		log.Println("incoming ...", msg.Text)

		if strings.Contains(msg.Text, yolow_command) {
			fBot.startLow(msg.Chat)
		} else if strings.Contains(msg.Text, yobg_command) {
			fBot.startQuickBg(msg.Chat)
		}
		if fBot.isRunning() {
			fBot.process.Run(msg)
		}
	}

}
