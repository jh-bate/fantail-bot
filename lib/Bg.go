package lib

import (
	"encoding/json"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type Bg struct {
	//messages   chan telebot.Message
	//Complete   chan bool
	Details    *Details
	SelectedBg string
	lang       struct {
		BgNow  []string      `json:"bgNow"`
		Review questionYesNo `json:"review"`
		Above  question      `json:"above"`
		In     question      `json:"in"`
		Below  question      `json:"below"`
		Thank  string        `json:"thank"`
	}
	Parts
}

func (this *Bg) loadLanguage() {

	encoded := `{
        "bgNow":  ["So your last blood sugar was ... "],
        "review": {
            "question": "Would you like to review your result?",
            "answers" : {
                "yes":"Sounds good",
                "no":"No thanks"
            }
        },
        "above": {
            "label":"Above what I would like",
            "question":"Do you have any idea why?",
            "children":[
                {
                    "label":"Yeah I think so",
                    "question":"Was it because ...",
                    "children":[
                        {
                            "label":"I guessed my BG",
                            "question":"Do you often guess your BG?",
                            "children":[
                                {
                                    "label":"Nope",
                                    "question":"",
                                    "children":null
                                },
                                {
                                    "label":" Yeah",
                                    "question":"",
                                    "children":null
                                }
                            ]
                        },
                        {
                            "label":"I guessed the carbs",
                            "question":"Do you often guess the amount of carbs?",
                            "children":[
                                {
                                    "label":"Nope",
                                    "question":"",
                                    "children":null
                                },
                                {
                                    "label":" Yeah",
                                    "question":"",
                                    "children":null
                                }
                            ]
                        },
                        {
                            "label":"I ate more the planned",
                            "question":"Is this common for you?",
                            "children":[
                                {
                                    "label":"Nope",
                                    "question":"",
                                    "children":null
                                },
                                {
                                    "label":" Yeah",
                                    "question":"",
                                    "children":null
                                }
                            ]
                        },
                        {
                            "label":"Just ate more",
                            "question":"",
                            "children":null
                        }
                    ]
                },
                {
                    "label":"No I don't",
                    "question":"Could it be because of any of the following reasons?",
                    "children":[
                        {
                            "label":"Guessed my BG",
                            "question":"",
                            "children": null
                        },
                        {
                            "label":"I guessed the carbs",
                            "question":"",
                            "children":null
                        },
                        {
                            "label":"I ate more the planned",
                            "question":"",
                            "children":null
                        },
                        {
                            "label":"Nope, none of the above",
                            "question":"",
                            "children":null
                        }
                    ]
                }
            ]
        },
        "in": {
            "label":"About right",
            "question":"Well done! So was that just as you planned :)",
            "children":[
                {
                    "label":"Yeah of course",
                    "question":"",
                    "children": null
                },
                {
                    "label":"Nope, but I will take it!",
                    "question":"",
                    "children": null
                }
            ]
        },
        "below": {
            "label":"Below what I would like",
            "question":"Do you have any idea why?",
            "children":[
                {
                    "label":"Yeah I think so",
                    "question":"Was it because ...",
                    "children":[
                        {
                            "label":"I didn't eat enough",
                            "question":"Do you know your insulin/carb ratio?",
                            "children":[
                                {
                                    "label":"Nope",
                                    "question":"",
                                    "children":null
                                },
                                {
                                    "label":" Yeah",
                                    "question":"",
                                    "children":null
                                }
                            ]
                        },
                        {
                            "label":"I had exercised",
                            "question":"Did you try and factor in the excercise?",
                            "children":[
                                {
                                    "label":"Nope",
                                    "question":"",
                                    "children":null
                                },
                                {
                                    "label":" Yeah",
                                    "question":"",
                                    "children":null
                                }
                            ]
                        },
                        {
                            "label":"I just stuffed up",
                            "question":"",
                            "children":null
                        }
                    ]
                },
                {
                    "label":"No I don't",
                    "question":"Could it be because of any of the following reasons?",
                    "children":[
                        {
                            "label":"I didn't eat enough",
                            "question":"",
                            "children": null
                        },
                        {
                            "label":"I had exercised",
                            "question":"",
                            "children":null
                        },
                        {
                            "label":"I don't know my insulin/carb ratio",
                            "question":"",
                            "children":null
                        },
                        {
                            "label":"Nope, none of the above",
                            "question":"",
                            "children":null
                        }
                    ]
                }
            ]
        },
        "thank":"Just wanted you to know - you rock! See ya."
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
		&Part{Func: bg.askBg, ToBeRun: true},
		&Part{Func: bg.askReview, ToBeRun: true},
		&Part{Func: bg.replyReview, ToBeRun: true},
	)
	return bg
}

/*func (this *Bg) AddMessages(msgs chan telebot.Message) {
	this.messages = msgs
}*/

/*func (this *Bg) Run(m telebot.Message) {
    done := make(chan bool)

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
}*/

func (this *Bg) Run(input <-chan telebot.Message) {

	for msg := range input {

		this.Details.User = msg.Chat

		log.Println("incoming", msg.Text)
		for i := range this.Parts {
			log.Println("checking ", i, "of", len(this.Parts), "still to run?", this.Parts[i].ToBeRun)
			if this.Parts[i].ToBeRun {
				log.Println("running ", i)
				this.Parts[i].ToBeRun = false
				this.Parts[i].Func(msg)
				log.Println("all done ", i)
				return
			}
		}
	}
}

/*func (this *Bg) CanRun() bool {
	return len(this.Parts) > 0
}*/

func (this *Bg) askBg(msg telebot.Message) {
	log.Println("ask the BG", this.lang.BgNow)
	log.Println("above?", this.lang.Above)
	this.Details.sendWithKeyboard(this.lang.BgNow[0], makeKeyBoard(this.lang.Above.Label, this.lang.In.Label, this.lang.Below.Label))
	return
}

func (this *Bg) askReview(msg telebot.Message) {
	this.SelectedBg = msg.Text
	this.Details.sendWithKeyboard(this.lang.Review.Question, makeKeyBoard(this.lang.Review.Yes, this.lang.Review.No))
	return
}

func (this *Bg) replyReview(msg telebot.Message) {
	switch {
	case msg.Text == this.lang.Review.Yes:
		this.doReview(this.SelectedBg)
		this.Parts = append(this.Parts, &Part{Func: this.onYa, ToBeRun: true})
		return
	case msg.Text == this.lang.Review.No:
		this.onYa(msg)
	}
	return
}

func (this *Bg) doReview(selectedAnswer string) {
	switch {
	case selectedAnswer == this.lang.Above.Label:
		this.addHighPath()
		return
	case selectedAnswer == this.lang.In.Label:
		this.addRightPath()
		return
	case selectedAnswer == this.lang.Below.Label:
		this.addLowPath()
		return
	}
	return
}

func (this *Bg) addHighPath() {
	for i := range this.lang.Above.Children {
		log.Println("adding ...", this.lang.Above.Children[i].Label)
		keys := []string{}
		for j := range this.lang.Above.Children[i].Children {
			log.Println("add key ...", this.lang.Above.Children[i].Children[j].Label)
			keys = append(keys, this.lang.Above.Children[i].Children[j].Label)
		}
		this.Details.sendWithKeyboard(this.lang.Above.Children[i].Question, makeKeyBoard(keys...))
	}
}

func (this *Bg) addLowPath() {
	for i := range this.lang.Below.Children {
		log.Println("adding ...", this.lang.Below.Children[i].Label)
		keys := []string{}
		for j := range this.lang.Below.Children[i].Children {
			log.Println("add key ...", this.lang.Below.Children[i].Children[j].Label)
			keys = append(keys, this.lang.Below.Children[i].Children[j].Label)
		}
		this.Details.sendWithKeyboard(this.lang.Below.Children[i].Question, makeKeyBoard(keys...))
	}
}

func (this *Bg) addRightPath() {
	for i := range this.lang.In.Children {
		log.Println("adding ...", this.lang.In.Children[i].Label)
		keys := []string{}
		for j := range this.lang.In.Children[i].Children {
			log.Println("add key ...", this.lang.In.Children[i].Children[j].Label)
			keys = append(keys, this.lang.In.Children[i].Children[j].Label)
		}
		this.Details.sendWithKeyboard(this.lang.In.Children[i].Question, makeKeyBoard(keys...))
	}
}

func (this *Bg) onYa(msg telebot.Message) {
	this.Details.sendWithKeyboard(this.lang.Thank, makeKeyBoard("See you!"))
	return
}
