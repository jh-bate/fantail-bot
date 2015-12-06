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
	Classify struct {
		classifier *bayesian.Classifier
	}

	Name bayesian.Class

	ClassificationWords []string

	ClassificationType struct {
		Name
		ClassificationWords
	}
)

func NewClassificationType(name Name, words ClassificationWords) *ClassificationType {
	return &ClassificationType{Name: name, ClassificationWords: words}
}

func (this Name) toBayesianClass() bayesian.Class {
	return bayesian.Class(this)
}

func NewClassification() *Classify {

	var LearnConfig struct {
		PostiveWords  ClassificationWords `json:"positive"`
		NegativeWords ClassificationWords `json:"negative"`
	}

	config.Load(&LearnConfig, "learn.json")

	pos := NewClassificationType(Postive, LearnConfig.PostiveWords)
	neg := NewClassificationType(Negative, LearnConfig.NegativeWords)
	classify := &Classify{
		classifier: bayesian.NewClassifier(pos.Name.toBayesianClass(), neg.Name.toBayesianClass()),
	}

	classify.classifier.Learn(pos.ClassificationWords, pos.toBayesianClass())
	classify.classifier.Learn(neg.ClassificationWords, neg.toBayesianClass())

	return classify
}

func (this *Classify) ArePositive(w ClassificationWords) bool {
	scores, likely, _ := this.classifier.LogScores(w)
	return scores[0] > scores[1] && likely == 0
}