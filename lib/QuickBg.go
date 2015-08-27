package lib

import (
	"encoding/json"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type QuickBg struct {
	Details *Details
	lang    struct {
		Comment  string `json:"comment"`
		Question string `json:"question"`
		Review   yesNo  `json:"review"`
		Above    option `json:"above"`
		In       option `json:"in"`
		Below    option `json:"below"`
		Thank    string `json:"thank"`
	}
	Parts
}

func (this *QuickBg) loadLanguage() {

	encoded := `{
        "comment": "OK lets save that blood sugar for you.",
        "question": "So your last blood sugar was ... ",
        "review" : {
        	"yes":"Sure thing!",
        	"no": "No thanks"
        },
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
        "thank":"You rock! Just had to say that before you go :)"
    }`

	err := json.Unmarshal([]byte(encoded), &this.lang)
	if err != nil {
		log.Panic("could not load BG language ", err.Error())
	}
}
func NewQuickBg(d *Details) *QuickBg {
	bg := &QuickBg{Details: d}
	bg.loadLanguage()
	bg.Parts = append(
		bg.Parts,
		&Part{Func: bg.selectBg, ToBeRun: true},
		&Part{Func: bg.questionReview, ToBeRun: true},
		&Part{Func: bg.answerReview, ToBeRun: true},
	)
	return bg
}

func (this *QuickBg) GetParts() Parts {
	return this.Parts
}

func (this *QuickBg) selectBg(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Question, makeKeyBoard(this.lang.Above.Text, this.lang.In.Text, this.lang.Below.Text))
	return
}

func (this *QuickBg) questionReview(msg telebot.Message) {
	this.Details.sendWithKeyboard(getLangText(this.lang.Above.FollowUpQuestion), makeKeyBoard(this.lang.Review.Yes, this.lang.Review.No))
	return
}

func (this *QuickBg) answerReview(msg telebot.Message) {
	switch {
	case msg.Text == this.lang.Review.Yes:
		this.Parts = append(
			this.Parts,
			&Part{Func: this.doReview, ToBeRun: true},
			&Part{Func: this.onYa, ToBeRun: true},
		)
	case msg.Text == this.lang.Review.No:
		this.Parts = append(
			this.Parts,
			&Part{Func: this.onYa, ToBeRun: true},
		)
	}
}

func (this *QuickBg) doReview(msg telebot.Message) {
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

func (this *QuickBg) onYa(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard("It does aye. See you!"))
	return
}
