package lib

import (
	"encoding/json"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type FoodWithFriends struct {
	Details        *Details
	SelectedAnswer string
	lang           struct {
		Comment           string   `json:"comment"`
		Instructions      []string `json:"instructions"`
		Understood        yesNo    `json:"understood"`
		Ready             yesNo    `json:"readyToGuess"`
		GuessInstructions []string `json:"guessInstructions"`
		Thank             string   `json:"thank"`
	}
	Parts
	snapResults struct {
		BgStart, BgFinish string
	}
}

//what are you eating
//take a pic
//ask others there estimated carb count

func (this *FoodWithFriends) loadLanguage() {

	encoded := `{
        "comment": "Right lets get this party started!",
        "instructions" : ["Take a clear pictue of what you are about to eat","Others will judge what you are about to eat from the photo and then guesstimate the carbs."],
        "understood" : {
        	"yes":"OK",
        	"no": "Tell me more ..."
        },
        "readyToGuess" : {
        	"yes":"Lets do it!",
        	"no": "Nope"
        },
        "guessInstructions" : ["Now you guess how many grams of carbs you think are in the food"],
        "thank":"You rock! Just had to say that before you go :)"
    }`

	err := json.Unmarshal([]byte(encoded), &this.lang)
	if err != nil {
		log.Panic("could not load BG language ", err.Error())
	}
}
func NewFoodWithFriends(d *Details) *FoodWithFriends {
	ff := &FoodWithFriends{Details: d}
	ff.loadLanguage()
	ff.Parts = append(
		ff.Parts,
		&Part{Func: ff.kickOff, ToBeRun: true},
		&Part{Func: ff.kickOffUnderstood, ToBeRun: true},
	)
	return ff
}

func (this *FoodWithFriends) Run(m telebot.Message) {
	for i := range this.Parts {
		if this.Parts[i].ToBeRun {
			this.Parts[i].ToBeRun = false
			this.Parts[i].Func(m)
			return
		}
	}
	return
}

func (this *FoodWithFriends) CanRun() bool {
	return len(this.Parts) > 0
}

func (this *FoodWithFriends) kickOff(msg telebot.Message) {
	this.Details.send(this.lang.Comment)
	this.Details.sendWithKeyboard(this.lang.Instructions[0], makeKeyBoard(this.lang.Understood.Yes, this.lang.Understood.No))
	return
}

func (this *FoodWithFriends) kickOffUnderstood(msg telebot.Message) {

	switch {
	case msg.Text == this.lang.Understood.Yes:
		this.Parts = append(this.Parts, &Part{Func: this.guessUnderstood, ToBeRun: true})
		this.Details.sendWithKeyboard(this.lang.GuessInstructions[0], makeKeyBoard(this.lang.Ready.Yes, this.lang.Ready.No))
		return
	case msg.Text == this.lang.Understood.No:
		this.Parts = append(this.Parts, &Part{Func: this.kickOffUnderstood, ToBeRun: true})
		this.Details.sendWithKeyboard(this.lang.Instructions[1], makeKeyBoard(this.lang.Understood.Yes, this.lang.Understood.No))
		return
	}
	return
}

func (this *FoodWithFriends) guessUnderstood(msg telebot.Message) {

	switch {
	case msg.Text == this.lang.Ready.Yes:
		this.Details.sendWithKeyboard(this.lang.GuessInstructions[0], makeKeyBoard(this.lang.Ready.Yes, this.lang.Ready.No))
		return
	case msg.Text == this.lang.Ready.No:
		this.Parts = append(this.Parts, &Part{Func: this.kickOffUnderstood, ToBeRun: true})
		this.Details.sendWithKeyboard(this.lang.Instructions[1], makeKeyBoard(this.lang.Understood.Yes, this.lang.Understood.No))
		return
	}
	return
}

func (this *FoodWithFriends) onYa(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard("See you!"))
	return
}
