package lib

import (
	"encoding/json"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type QuickBg struct {
	Details        *Details
	SelectedAnswer string
	lang           struct {
		Comment     string `json:"comment"`
		Question    string `json:"question"`
		Review      string `json:"review"`
		ReviewYesNo yesNo  `json:"reviewYesNo"`
		Above       option `json:"above"`
		In          option `json:"in"`
		Below       option `json:"below"`
		Thank       string `json:"thank"`
	}
	Parts
}

func (this *QuickBg) loadLanguage() {

	encoded := `{
        "comment": "OK lets save that blood sugar for you.",
        "question": "So your last blood sugar was ... ",
        "review": "Would you like to review your result?",
        "reviewYesNo" : {
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
		&Part{Func: bg.askBg, ToBeRun: true},
		&Part{Func: bg.askReview, ToBeRun: true},
		&Part{Func: bg.replyReview, ToBeRun: true},
	)
	return bg
}

func (this *QuickBg) Run(m telebot.Message) {
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

func (this *QuickBg) CanRun() bool {
	return len(this.Parts) > 0
}

func (this *QuickBg) askBg(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Question, makeKeyBoard(this.lang.Above.Text, this.lang.In.Text, this.lang.Below.Text))
	return
}

func (this *QuickBg) askReview(msg telebot.Message) {
	this.SelectedAnswer = msg.Text
	this.Details.sendWithKeyboard(this.lang.Review, makeKeyBoard(this.lang.ReviewYesNo.Yes, this.lang.ReviewYesNo.No))
	return
}

func (this *QuickBg) replyReview(msg telebot.Message) {
	log.Println("did you want to review?", msg.Text)
	switch {
	case msg.Text == this.lang.ReviewYesNo.Yes:
		this.Parts = append(
			this.Parts,
			&Part{Func: this.doReview, ToBeRun: true},
			&Part{Func: this.onYa, ToBeRun: true},
		)
	case msg.Text == this.lang.ReviewYesNo.No:
		this.Parts = append(
			this.Parts,
			&Part{Func: this.onYa, ToBeRun: true},
		)
	}
	return
}

func (this *QuickBg) doReview(msg telebot.Message) {
	log.Println("BG", this.SelectedAnswer)
	switch {
	case this.SelectedAnswer == this.lang.Above.Text:
		this.Details.send(getLangText(this.lang.Above.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.Above.FollowUpQuestion), makeKeyBoard("Sure has", "Nope"))
		return
	case this.SelectedAnswer == this.lang.In.Text:
		this.Details.send(getLangText(this.lang.In.Feedback))
		this.Details.sendWithKeyboard(getLangText(this.lang.In.FollowUpQuestion), makeKeyBoard("Totally!", "Not so sure"))
		return
	case this.SelectedAnswer == this.lang.Below.Text:
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
