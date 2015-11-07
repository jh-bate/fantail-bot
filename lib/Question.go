package lib

import "log"

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

func (this Questions) next(prevAnswer string) (*Question, bool) {
	for i := range this {
		for a := range this[i].RelatesTo.Answers {
			log.Println(this[i].RelatesTo.Answers[a], "matches", prevAnswer)
			if this[i].RelatesTo.Answers[a] == prevAnswer {
				return this[i], this[i].RelatesTo.Save
			}
		}
	}
	return nil, false
}

func (this Questions) nextFrom(prevAnswers ...string) (*Question, bool) {
	for i := range prevAnswers {
		nxt, sv := this.next(prevAnswers[i])
		if nxt != nil {
			return nxt, sv
		}
	}
	return nil, false
}

func (this *Question) makeKeyboard() Keyboard {
	keyboard := Keyboard{}
	for i := range this.PossibleAnswers {
		keyboard = append(keyboard, []string{this.PossibleAnswers[i]})
	}
	return keyboard
}
