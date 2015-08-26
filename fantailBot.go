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

	typing_action = "typing"
)

type (
	fantailBot struct {
		bot *telebot.Bot
		*running
	}

	running struct {
		Parts lib.Parts
		Name  string
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

func (this *fantailBot) setRunning(p lib.Parts, n string) {
	this.running = &running{Parts: p, Name: n}
	return
}

func (this *fantailBot) isRunning(name string) bool {
	return this.running != nil && this.running.Name == name
}

func main() {

	fBot := newFantailBot()
	messages := make(chan telebot.Message)
	fBot.bot.Listen(messages, 1*time.Second)

	for msg := range messages {

		log.Println("incoming ...", msg.Text)

		if strings.Contains(msg.Text, yobg_command) || fBot.isRunning(yobg_command) {
			if fBot.isRunning(yobg_command) == false {
				fBot.setRunning(
					lib.NewBg(&lib.Details{Bot: fBot.bot, User: msg.Chat}).GetParts(),
					yobg_command,
				)
			}
			lib.Run(msg, fBot.running.Parts)
		} else if strings.Contains(msg.Text, yolow_command) || fBot.isRunning(yolow_command) {
			if fBot.isRunning(yolow_command) == false {
				fBot.setRunning(
					lib.NewLow(&lib.Details{Bot: fBot.bot, User: msg.Chat}).GetParts(),
					yolow_command,
				)
			}
			lib.Run(msg, fBot.running.Parts)
		} else {
			fBot.showOptions(msg.Chat)
		}
	}

}
