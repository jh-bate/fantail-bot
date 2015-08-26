package lib

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type Low struct {
	Details *Details
	lang    struct {
		Comment  string `json:"comment"`
		Question string `json:"question"`
		Good     option `json:"good"`
		NotGood  option `json:"notGood"`
		Other    option `json:"other"`
		Thank    string `json:"thank"`
	}
	Parts
}

func (this *Low) loadLanguage() {

	encoded := `{
	  "comment": "Stink! We hope you are back on track now.",
	  "question": "So how do you feel you coped with the low?",
	  "good":{
	    "text":"OK",
	    "feedback":["Well done! Lows are a pain to great to hear that you knocked it on the head"],
	    "followUp" :["Do you have any idea why you went low?"]
	  },
	  "notGood":{
	    "text":"Not that well",
	    "feedback":["Sometimes the temptation is to rush things ... it pays to slow down sometimes"],
	    "followUp" :["Do you have an idea of how you would better deal with it next time?","Do you have any idea why you went low?"]
	  },
	  "other":{
	    "text":"You know how it goes",
	    "feedback":["Yeap, we sure do!","Just remember sharing is caring :)"],
	    "followUp" :["Maybe next time?"]
	  },
	  "thank":"Lets plan on no more of those to deal with for a while!"
	}`

	err := json.Unmarshal([]byte(encoded), &this.lang)
	if err != nil {
		log.Panic("could not load LOW language ", err.Error())
	}
}

func NewLow(d *Details) *Low {
	low := &Low{Details: d}
	low.loadLanguage()
	low.Parts = append(
		low.Parts,
		&Part{Func: low.partOne, ToBeRun: true},
		&Part{Func: low.partTwo, ToBeRun: true},
		&Part{Func: low.partThree, ToBeRun: true},
	)
	return low
}

func (this *Low) GetParts() Parts {
	return this.Parts
}

func (this *Low) partOne(msg telebot.Message) {
	this.Details.send(fmt.Sprintf("Hey %s", msg.Chat.FirstName))
	this.Details.send(this.lang.Comment)
	this.Details.sendWithKeyboard(this.lang.Question, makeKeyBoard(this.lang.Good.Text, this.lang.NotGood.Text, this.lang.Other.Text))
	return
}

func (this *Low) partTwo(msg telebot.Message) {
	switch {
	case msg.Text == this.lang.Good.Text:
		this.Details.send(getLangText(this.lang.Good.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.Good.FollowUpQuestion), makeKeyBoard(yes_text, no_text))
		return
	case msg.Text == this.lang.NotGood.Text:
		this.Details.send(getLangText(this.lang.NotGood.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.NotGood.FollowUpQuestion), makeKeyBoard(yes_text, no_text))
		return
	case msg.Text == this.lang.Other.Text:
		this.Details.send(getLangText(this.lang.Other.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.Other.FollowUpQuestion), makeKeyBoard(yes_text, no_text))
		return
	}
	return
}

func (this *Low) partThree(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard(bye_text))
	return
}
