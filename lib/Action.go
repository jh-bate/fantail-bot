package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

const (
	start_action, help_action = "/start", "/help"

	chat_action                       = "/chat"
	say_action, say_action_hint       = "/say", "/say [what you want to say]"
	recap_action, recap_action_hint   = "/recap", "/recap <days> <what>"
	remind_action, remind_action_hint = "/remind in <days> to <msg>", "/remind"

	//for any submisson
	default_action = "default"
)

type (
	action struct {
		typeOf  string
		context string
	}
)

var (
	actionTypes = []string{chat_action, say_action, recap_action, remind_action}
)

func newAction(msg telebot.Message) *action {

	if msg.Sticker.Exists() {
		return &action{typeOf: chat_action}
	}

	words := strings.Fields(msg.Text)

	if strings.Contains(words[0], "/") {
		for i := range actionTypes {
			if strings.ToLower(actionTypes[i]) == strings.ToLower(words[0]) {
				return &action{typeOf: words[0], context: msg.Text}
			}
		}
	}

	return &action{typeOf: "", context: msg.Text}
}

func (this action) typeIsSet() bool {
	return this.typeOf != ""
}

func (this action) setTypeMatches(types ...string) bool {
	for i := range types {
		if types[i] == this.typeOf {
			return true
		}
	}
	return false
}

func (this action) isHelp() bool {
	return this.typeOf == start_action || this.typeOf == help_action
}

func (this action) info() string {
	if this.isHelp() {
		return fmt.Sprintf(" %s %s %s ",
			chat_action+" - to have a *quick chat* about what your upto \n\n",
			say_action_hint+" - to say *anything* thats on your mind \n\n",
			recap_action_hint+" - to go over what has been happening \n\n",
		)
	}
	return ""
}

func (this action) loadQuestions(q Questions) {

	if this.typeIsSet() {

		name := strings.SplitAfter(this.typeOf, "/")[1]
		if this.hasSubmisson() {
			name = default_action
		}

		file, err := os.Open(fmt.Sprintf("./config/%s.json", name))
		if err != nil {
			log.Panic("could not load QandA language file ", err.Error())
		}
		err = json.NewDecoder(file).Decode(&q)
		if err != nil {
			log.Panic("could not decode QandA ", err.Error())
		}
	}
	return
}

func (this action) hasSubmisson() bool {

	words := strings.Fields(this.context)
	if len(words) > 1 {
		return true
	}

	return false
}

func (this action) submissionDays() int {
	const cmd_pos, time_pos = 0, 1
	words := strings.Fields(this.context)

	if len(words) == 1 {
		//dafault if zero
		return 0
	}

	days, err := strconv.Atoi(words[time_pos])

	if err != nil {
		log.Println("error getting number of days", err.Error())
	}

	log.Println("days", days)

	return days
}
