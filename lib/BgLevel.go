package lib

import (
	"encoding/json"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type BgLevel struct {
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
                    "question":"Do you know why it was above?",
                    "children":[
                        {
                            "label":"Yeah I think so",
                            "question":"Was it above because ...",
                            "children":[
                                {
                                    "label":"I guessed my BG before eating",
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
                                    "label":"I didn't know how many carbs I actually ate",
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
                                    "label":"A bad site or injection",
                                    "question":"Do you know how to deal with this?",
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
                                    "label":"I am sick / I think I am gettig sick",
                                    "question":"Do you have a plan in place for when your sick?",
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
                            "label":"Yeah of course :)",
                            "question":"",
                            "children": null
                        },
                        {
                            "label":"Nope, but I will take the win!",
                            "question":"",
                            "children": null
                        }
                    ]
                },
                {
                    "label":"Below what I would like",
                    "question":"Do you know why it was below?",
                    "children":[
                        {
                            "label":"Yeah",
                            "question":"Was it below because ...",
                            "children":[
                                {
                                    "label":"I got my carb-to-insulin ratio wrong",
                                    "question":"Do you know your carb-to-insulin ratio?",
                                    "children":[
                                        {
                                            "label":"Yeah",
                                            "question":"",
                                            "children":null
                                        },
                                        {
                                            "label":"No",
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
                            "label":"Nah",
                            "question":"Could it be because of any of the following reasons?",
                            "children":[
                                {
                                    "label":"I had exercised",
                                    "question":"",
                                    "children":null
                                },
                                {
                                    "label":"Your unsure of the carb-to-insulin ratio?",
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

/*func (this *BgLevel) loadConfig() {
	config, err := ioutil.ReadFile("./bgConfig.json")
	if err != nil {
		log.Fatal("Loading BG config", err.Error())
	}
	err = json.Unmarshal(config, &this.lang)
	if err != nil {
		log.Panic("Unmarshaling BG config", err.Error())
	}
	return
}*/

func NewBgLevel(d *Details) *BgLevel {
	bg := &BgLevel{Details: d}
	bg.loadLanguage()
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
	nextQ := this.lang.BgNow.find(msg.Text)
	if nextQ == nil {
		this.Details.send(this.lang.Thank)
		log.Println("all done now")
		return
	}
	log.Println("asking ...", nextQ.Question, "labeled:", nextQ.Label)
	this.Details.sendWithKeyboard(nextQ.Question, nextQ.keyboard())
	return
}
