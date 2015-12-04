package lib

import (
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jbrukh/bayesian"
)

const (
	Postive  Name = "Postive"
	Negative Name = "Negative"
)

type (
	Learn struct {
		classifier *bayesian.Classifier
	}

	Name bayesian.Class

	LearningType struct {
		Name
		Words []string `json:"words"`
	}
)

func NewLearningType(name Name) *LearningType {
	return &LearningType{Name: name}
}

func (this Name) toBayesianClass() bayesian.Class {
	return bayesian.Class(this)
}

func (this *LearningType) loadWords() {

	if this.Name == Postive {
		this.Words = []string{
			"happy",
			"great",
			"won",
			"good",
		}
	} else if this.Name == Negative {
		this.Words = []string{
			"low",
			"high",
			"depressed",
			"over it",
			"sick",
		}
	}
	return
}

func NewLearner() *Learn {

	pos := NewLearningType(Postive)
	pos.loadWords()
	neg := NewLearningType(Negative)
	neg.loadWords()
	mc := &Learn{
		classifier: bayesian.NewClassifier(pos.Name.toBayesianClass(), neg.Name.toBayesianClass()),
	}

	mc.classifier.Learn(pos.Words, pos.toBayesianClass())
	mc.classifier.Learn(neg.Words, neg.toBayesianClass())

	return mc
}

func (this *Learn) isPositive(n Notes) bool {

	var noteText []string

	for i := range n {
		noteText = append(noteText, strings.Fields(n[i].Text)...)
	}

	scores, likely, _ := this.classifier.LogScores(noteText)

	return scores[0] > scores[1] && likely == 0

}
