package lib

import (
	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jbrukh/bayesian"

	"github.com/jh-bate/fantail-bot/lib/config"
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

	Words []string

	LearningType struct {
		Name
		Words
	}
)

func NewLearningType(name Name, words Words) *LearningType {
	return &LearningType{Name: name, Words: words}
}

func (this Name) toBayesianClass() bayesian.Class {
	return bayesian.Class(this)
}

func NewLearner() *Learn {

	var LearnConfig struct {
		PostiveWords  Words `json:"positive"`
		NegativeWords Words `json:"negative"`
	}

	config.Load(&LearnConfig, "learn.json")

	pos := NewLearningType(Postive, LearnConfig.PostiveWords)
	neg := NewLearningType(Negative, LearnConfig.NegativeWords)
	mc := &Learn{
		classifier: bayesian.NewClassifier(pos.Name.toBayesianClass(), neg.Name.toBayesianClass()),
	}

	mc.classifier.Learn(pos.Words, pos.toBayesianClass())
	mc.classifier.Learn(neg.Words, neg.toBayesianClass())

	return mc
}

func (this *Learn) ArePositive(w Words) bool {
	scores, likely, _ := this.classifier.LogScores(w)
	return scores[0] > scores[1] && likely == 0
}
