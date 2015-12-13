package question

import (
	"log"
	"strings"
)

type (
	Question struct {
		RelatesTo struct {
			Answers []string `json:"answers"`
			Save    bool     `json:"save"`
			SaveTag string   `json:"saveTag"`
		} `json:"relatesTo"`
		Context         []string `json:"context"`
		QuestionText    string   `json:"question"`
		PossibleAnswers []string `json:"answers"`
	}

	Questions []*Question
)

func (this Questions) First() *Question {
	if len(this) > 0 {
		return this[0]
	}
	return nil
}

func (this Questions) Next(prevAnswer string) (*Question, bool) {
	for i := range this {
		for a := range this[i].RelatesTo.Answers {
			if strings.EqualFold(this[i].RelatesTo.Answers[a], prevAnswer) {
				return this[i], this[i].RelatesTo.Save
			}
		}
	}
	return nil, false
}

func (this Questions) NextFrom(prevAnswers ...string) (*Question, bool) {
	for i := range prevAnswers {
		if nxt, sv := this.Next(prevAnswers[i]); nxt != nil {
			log.Println("got it from sticker ...")
			return nxt, sv
		}
	}
	return nil, false
}

func (this *Question) MakeKeyboard() [][]string {
	keyboard := [][]string{}
	for i := range this.PossibleAnswers {
		keyboard = append(keyboard, []string{this.PossibleAnswers[i]})
	}
	return keyboard
}
