package lib

import (
	"encoding/json"
	"log"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type Food struct {
	Details        *Details
	SelectedAnswer string
	lang           struct {
		Comment   string `json:"comment"`
		Question  string `json:"question"`
		Meal      string `json:"mean"`
		Snack     string `json:"snack"`
		Beverage  string `json:"beverage"`
		PlaySnap  string `json:"playSnap"`
		SnapYesNo yesNo  `json:"snapYesNo"`
		Snapped   yesNo  `json:"snapped"`
		Thank     string `json:"thank"`
	}
	Parts
	snapResults struct {
		BgStart, BgFinish string
	}
}

//what are you eating
//lets play a game of snap
//bg before = bg after = snap

func (this *Food) loadLanguage() {

	encoded := `{
        "comment": "OK lets save that blood sugar for you.",
        "question": "So what are you about to have  ... ",
        "meal": "A Meal",
        "snack": "A Snack",
        "beverage": "A Drink",
        "playSnap": "Would you like to play a game of snap?/n If you BG before equals BG after then its SNAP!",
        "snapYesNo" : {
        	"yes":"Lets do it!",
        	"no": "No thanks"
        },
        "snapped" : {
        	"yes":"SNAP!!",
        	"no": "Next time!"
        },
        "thank":"You rock! Just had to say that before you go :)"
    }`

	err := json.Unmarshal([]byte(encoded), &this.lang)
	if err != nil {
		log.Panic("could not load BG language ", err.Error())
	}
}
func NewFood(d *Details) *Food {
	bg := &Food{Details: d}
	bg.loadLanguage()
	bg.Parts = append(
		bg.Parts,
		&Part{Func: bg.whatFood, ToBeRun: true},
		&Part{Func: bg.askToPlaySnap, ToBeRun: true},
		&Part{Func: bg.areWePalyingSnap, ToBeRun: true},
	)
	return bg
}

func (this *Food) Run(m telebot.Message) {
	for i := range this.Parts {
		if this.Parts[i].ToBeRun {
			this.Parts[i].ToBeRun = false
			this.Parts[i].Func(m)
			return
		}
	}
	return
}

func (this *Food) CanRun() bool {
	return len(this.Parts) > 0
}

func (this *Food) whatFood(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Question, makeKeyBoard(this.lang.Beverage, this.lang.Snack, this.lang.Meal))
	return
}

func (this *Food) askToPlaySnap(msg telebot.Message) {
	this.SelectedAnswer = msg.Text
	this.Details.sendWithKeyboard(this.lang.PlaySnap, makeKeyBoard(this.lang.SnapYesNo.Yes, this.lang.SnapYesNo.No))
	return
}

func (this *Food) areWePalyingSnap(msg telebot.Message) {
	switch {
	case msg.Text == this.lang.SnapYesNo.Yes:
		bg := NewQuickBg(this.Details)
		bg.askBg(msg)
		this.Parts = append(
			this.Parts,
			&Part{Func: this.startSnap, ToBeRun: true},
			&Part{Func: this.finishSnap, ToBeRun: true},
			&Part{Func: this.onYa, ToBeRun: true},
		)
		return
	case msg.Text == this.lang.SnapYesNo.No:
		this.onYa(msg)
	}
	return
}

func (this *Food) startSnap(msg telebot.Message) {
	this.snapResults.BgStart = msg.Text
	time.Sleep(30 * time.Second)
	return
}

func (this *Food) finishSnap(msg telebot.Message) {
	this.snapResults.BgFinish = msg.Text

	if this.snapResults.BgStart == this.snapResults.BgFinish {
		this.Details.sendWithKeyboard(this.lang.Snapped.Yes, makeKeyBoard("Woot!"))
	} else {
		this.Details.sendWithKeyboard(this.lang.Snapped.No, makeKeyBoard("Next time!"))
	}

	return
}

func (this *Food) onYa(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard("See you!"))
	return
}
