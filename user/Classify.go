package user

import (
	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/jbrukh/bayesian"
	"github.com/jh-bate/fantail-bot/config"
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

	ClassificationType struct {
		Name
		Words []string
	}
)

func NewClassificationType(name Name, words []string) *ClassificationType {
	return &ClassificationType{Name: name, Words: words}
}

func (this Name) toBayesianClass() bayesian.Class {
	return bayesian.Class(this)
}

func NewClassification() *Classify {

	var ClassifyConfig struct {
		PostiveWords  []string `json:"positive"`
		NegativeWords []string `json:"negative"`
	}

	config.Load(&ClassifyConfig, "classify.json")

	pos := NewClassificationType(Postive, ClassifyConfig.PostiveWords)
	neg := NewClassificationType(Negative, ClassifyConfig.NegativeWords)
	classify := &Classify{
		classifier: bayesian.NewClassifier(pos.Name.toBayesianClass(), neg.Name.toBayesianClass()),
	}

	classify.classifier.Learn(pos.Words, pos.toBayesianClass())
	classify.classifier.Learn(neg.Words, neg.toBayesianClass())

	return classify
}

func (this *Classify) ArePositive(w []string) bool {
	scores, likely, _ := this.classifier.LogScores(w)
	return scores[0] > scores[1] && likely == 0
}
