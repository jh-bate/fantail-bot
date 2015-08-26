package lib

import (
	"encoding/json"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

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

const bg_config_name = "bg.json"

func NewBg(d *Details) *Bg {
	bg := &Bg{Details: d}
	bg.loadLanguage()
	bg.Parts = append(
		bg.Parts,
		&Part{Func: bg.partOne, ToBeRun: true},
		&Part{Func: bg.partTwo, ToBeRun: true},
		&Part{Func: bg.partThree, ToBeRun: true},
	)
	return bg
}

func (this *Bg) loadLanguage() {

	encoded := `{
	  "comment": "Cool lets get that done for you then.",
	  "question": "So your last bloodsugar was ... ",
		"above": {
	    "text": "Above what I would like",
	    "feedback":["Hydrate!","Say bye to high!","Just remember it happens to the best of us"],
	     "followUp" :["Has it been high for a while?"]
	  },
		"in": {
	    "text": "About right",
	    "feedback":["Awesome work!!","Its never as easy as its made out aye :)","How does it feel to be perfect :)"],
	    "followUp" :["Did you feel you could do this again and again?"]
	  },
		"below":{
	    "text": "Below what I would like",
	    "feedback":["Lets say no to low!","Damn lows","Hope you keep your low supplies stocked up"],
	    "followUp" :["Do you have any idea why you went low?"]
	  },
	  "thank":"Thanks for that - it all counts"
	}`

	err := json.Unmarshal([]byte(encoded), &this.lang)
	if err != nil {
		log.Panic("could not load BG language ", err.Error())
	}
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
