package lib

import (
	"encoding/json"
	"log"
	"os"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type QTree struct {
	Details *Details
	lang    struct {
		QTree questionTree `json:"QTree"`
		Thank string       `json:"thank"`
	}
}

func NewQTree(d *Details) *QTree {
	bg := &QTree{Details: d}
	bg.loadLanguage()
	return bg
}

func (this *QTree) Run(input <-chan telebot.Message) {
	for msg := range input {
		this.Details.User = msg.Chat
		log.Println("in", msg.Text)
		this.ask(msg)
	}
}

func (this *QTree) loadLanguage() {

	file, err := os.Open("./config/why_high.json")
	if err != nil {
		log.Panic("could not load QTree language file ", err.Error())
	}
	err = json.NewDecoder(file).Decode(&this.lang)
	if err != nil {
		log.Panic("could not decode QTree ", err.Error())
	}
}

func (q *questionTree) hasChildren() bool {
	return q.Children != nil && len(q.Children) > 0
}

func (q *questionTree) find(label string) *questionTree {
	if q.hasChildren() {

		if q.Label == label {
			log.Println("at the top so return this one's children")
			return q
		}
		for i := range q.Children {
			if q.Children[i].Label == label {
				return q.Children[i]
			} else if q.Children[i].hasChildren() {
				match := q.Children[i].find(label)
				if match != nil {
					return match
				}
			}
		}
	}
	log.Println("nothing else found")
	return nil
}

func (this *QTree) makeKeyboard(q *questionTree) Keyboard {
	keyboard := Keyboard{}
	for i := range q.Children {
		keyboard = append(keyboard, []string{q.Children[i].Label})
	}
	return keyboard
}

func (this *QTree) askQuestion(q *questionTree) {
	//log.Println("asking ...", q.Label)
	for i := range q.Questions {
		if i+1 != len(q.Questions) {
			//log.Println("snd msg")
			this.Details.send(q.Questions[i])
		} else {
			//log.Println("snd msg w kb")
			this.Details.sendWithKeyboard(q.Questions[i], this.makeKeyboard(q))
		}
	}
	return
}

func (this *QTree) ask(msg telebot.Message) {
	//log.Println("answer was", msg.Text)
	nextQ := this.lang.QTree.find(msg.Text)
	if nextQ == nil {
		this.Details.send(this.lang.Thank)
		log.Println("all done now")
		return
	}

	this.askQuestion(nextQ)
	return
}
