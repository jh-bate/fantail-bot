package lib

import (
	"encoding/json"
	"fmt"
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

func (this *Bg) loadLanguage() {

	encoded := `{
        "comment": "OK lets save that blood sugar for you.",
        "question": "So your last blood sugar was ... ",
        "above": {
            "text": "Above what I would like",
            "feedback":["Remember to keep your fuilds up","Just remember it happens to the best of us"],
            "followUp" :["Has it been high for a while?"]
        },
        "in": {
            "text": "About right",
            "feedback":["Awesome work!!","Its never as easy as its made out aye :)","How does it feel to be perfect :)"],
            "followUp" :["Did you feel you could do this again and again?"]
        },
        "below":{
            "text": "Below what I would like",
            "feedback":["Damn lows, they always happen at the wrost time.","I Hope you keep your low supplies stocked up!!"],
            "followUp" :["Do you have any idea why you went low?"]
        },
        "thank":"Thanks for that - it all counts"
    }`

	err := json.Unmarshal([]byte(encoded), &this.lang)
	if err != nil {
		log.Panic("could not load BG language ", err.Error())
	}
}
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

func (this *Bg) Run(m telebot.Message) {
	for i := range this.Parts {
		log.Println("checking ", i, "of", len(this.Parts), "still to run?", this.Parts[i].ToBeRun)
		if this.Parts[i].ToBeRun {
			log.Println("running ", i)
			this.Parts[i].ToBeRun = false
			this.Parts[i].Func(m)
			log.Println("all done ", i)
			return
		}
	}
	return
}

func (this *Bg) CanRun() bool {
	return len(this.Parts) > 0
}

func (this *Bg) partOne(msg telebot.Message) {
	this.Details.send(fmt.Sprintf("Hey %s", msg.Chat.FirstName))
	this.Details.send(this.lang.Comment)
	this.Details.sendWithKeyboard(this.lang.Question, makeKeyBoard(this.lang.Above.Text, this.lang.In.Text, this.lang.Below.Text))
	return
}

func (this *Bg) partTwo(msg telebot.Message) {
	switch {
	case msg.Text == this.lang.Above.Text:
		this.Details.send(getLangText(this.lang.Above.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.Above.FollowUpQuestion), makeKeyBoard("Sure has", "Nope"))
		return
	case msg.Text == this.lang.In.Text:
		this.Details.send(getLangText(this.lang.In.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.In.FollowUpQuestion), makeKeyBoard("Totally!", "Not so sure"))
		return
	case msg.Text == this.lang.Below.Text:
		this.Details.send(getLangText(this.lang.Below.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.Below.FollowUpQuestion), makeKeyBoard("Yeah I have a hunch", "No, I just don't get it"))
		return
	}
	return
}

func (this *Bg) partThree(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard("It does aye. See you!"))
	return
}
