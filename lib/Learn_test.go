package lib

import (
	"strings"
	"testing"
)

func test_makeHappyWords() Words {
	w := Words{}
	w = append(w, strings.Fields("/say really happy")...)
	w = append(w, strings.Fields("/say hello all good")...)
	w = append(w, strings.Fields("things very good, went really well today")...)
	w = append(w, strings.Fields("/say bad day today, too many lows")...)
	return w
}

func test_makeNeutralWords() Words {
	w := Words{}
	w = append(w, strings.Fields("/say low happy")...)
	w = append(w, strings.Fields("/say bad happy")...)
	w = append(w, strings.Fields("all going well bad")...)
	w = append(w, strings.Fields("/say help good")...)
	w = append(w, strings.Fields("/say low great")...)
	return w

}

func test_makeUnhappyWords() Words {
	w := Words{}
	w = append(w, strings.Fields("/say help")...)
	w = append(w, strings.Fields("/say more highs, sick of it!")...)
	w = append(w, strings.Fields("things went really well today")...)
	w = append(w, strings.Fields("/say all good help")...)
	w = append(w, strings.Fields("/say low again!!")...)
	return w
}

func TestLearn_Postive(t *testing.T) {

	learn := NewLearner()

	if !learn.ArePositive(test_makeHappyWords()) {
		t.Error("this should have been recorded as positive")
	}

}

func TestLearn_Postive_WhenBalanced(t *testing.T) {

	learn := NewLearner()

	if !learn.ArePositive(test_makeNeutralWords()) {
		t.Error("this should have been recorded as positive when notes are neutral")
	}

}

func TestLearn_Negative(t *testing.T) {

	learn := NewLearner()

	if learn.ArePositive(test_makeUnhappyWords()) {
		t.Error("this should have been recorded as negative")
	}

}
