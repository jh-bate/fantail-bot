package lib

import "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"

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

func NewLow(d *Details) *Low {
	low := &Low{Details: d}
	loadLanguage(low.lang)
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
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard(yes_text, no_text))
	return
}
