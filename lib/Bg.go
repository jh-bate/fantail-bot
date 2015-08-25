package lib

import "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"

type Bg struct {
	Details *Details
	lang    struct {
		Comment  string `json:"comment"`
		Question string `json:"question"`
		Above    option `json:"above"`
		In       option `json:"in"`
		Below    option `json:"below"`
		Thank    string `json:"thank"`
	}
	Parts
}

func NewBg(d *Details) *Bg {
	bg := &Bg{Details: d}
	loadLanguage(bg.lang)
	bg.Parts = append(
		bg.Parts,
		&Part{Func: bg.partOne, ToBeRun: true},
		&Part{Func: bg.partTwo, ToBeRun: true},
		&Part{Func: bg.partThree, ToBeRun: true},
	)
	return bg
}

func (this *Bg) GetParts() Parts {
	return this.Parts
}

func (this *Bg) partOne(msg telebot.Message) {
	this.Details.send(this.lang.Comment)
	this.Details.sendWithKeyboard(this.lang.Question, makeKeyBoard(this.lang.Above.Text, this.lang.In.Text, this.lang.Below.Text))
	return
}

func (this *Bg) partTwo(msg telebot.Message) {
	switch {
	case msg.Text == this.lang.Above.Text:
		this.Details.send(getLangText(this.lang.Above.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.Above.FollowUpQuestion), makeKeyBoard(yes_text, no_text))
		return
	case msg.Text == this.lang.In.Text:
		this.Details.send(getLangText(this.lang.In.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.In.FollowUpQuestion), makeKeyBoard(yes_text, no_text))
		return
	case msg.Text == this.lang.Below.Text:
		this.Details.send(getLangText(this.lang.Below.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.Below.FollowUpQuestion), makeKeyBoard(yes_text, no_text))
		return
	}
	return
}

func (this *Bg) partThree(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard(yes_text, no_text))
	return
}
