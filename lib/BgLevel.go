package lib

import (
	"encoding/json"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type BgLevel struct {
	//Done    chan bool
	Details *Details
	lang    struct {
		BgNow question `json:"bgNow"`
		Thank string   `json:"thank"`
	}
}

func (this *BgLevel) loadLanguage() {

	encoded := `{
        "bgNow":  {
            "label":"/bg",
            "question":"So your last blood sugar was ...",
            "children":[
                {
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
                {
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
                {
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
                }
            ]
        },
        "thank":"Just wanted you to know - you rock! See ya."
    }`

	err := json.Unmarshal([]byte(encoded), &this.lang)
	if err != nil {
		log.Panic("could not load BgLevel language ", err.Error())
	}
}

func NewBgLevel(d *Details) *BgLevel {
	bg := &BgLevel{Details: d}
	bg.loadLanguage()
	//bg.Done = make(chan bool)
	return bg
}

func (this *BgLevel) Run(input <-chan telebot.Message) {
	for msg := range input {
		this.Details.User = msg.Chat
		this.ask(msg)
	}
}

func (this *BgLevel) ask(msg telebot.Message) {
	log.Println("answer was", msg.Text)
	nextQ := this.lang.BgNow.findChild(msg.Text)
	if nextQ == nil {
		this.Details.send(this.lang.Thank)
		log.Println("all done now")
		//close(this.Done)
		return
	}
	log.Println("asking ...", nextQ.Question, "labeled:", nextQ.Label)
	this.Details.sendWithKeyboard(nextQ.Question, nextQ.keyboard())
	return
}
