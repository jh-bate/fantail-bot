package lib

import "testing"

func test_makeHappyNotes() Notes {
	return Notes{
		NewNote(newMsg("/say really happy")),
		NewNote(newMsg("/say hello")),
		NewNote(newMsg("all going well")),
		NewNote(newMsg("/say need help to figure stuff out"), help_tag),
		NewNote(newMsg("/say bad day today, too many lows")),
	}
}

func test_makeNeutralNotes() Notes {
	return Notes{
		NewNote(newMsg("/say low happy")),
		NewNote(newMsg("/say bad happy")),
		NewNote(newMsg("all going well bad")),
		NewNote(newMsg("/say help good"), help_tag),
		NewNote(newMsg("/say low great")),
	}
}

func test_makeUnhappyNotes() Notes {
	return Notes{
		NewNote(newMsg("/say help")),
		NewNote(newMsg("/say more highs, sick of it!")),
		NewNote(newMsg("things went really well today")),
		NewNote(newMsg("/say all good"), help_tag),
		NewNote(newMsg("/say low again!!")),
	}
}

func TestLearn_Postive(t *testing.T) {

	learn := NewLearner()

	if !learn.isPositive(test_makeHappyNotes()) {
		t.Error("this should have been recoreded as positive")
	}

}

func TestLearn_Postive_WhenBalanced(t *testing.T) {

	learn := NewLearner()

	if !learn.isPositive(test_makeNeutralNotes()) {
		t.Error("this should have been recoreded as positive when notes are nuteral")
	}

}

func TestLearn_Negative(t *testing.T) {

	learn := NewLearner()

	if learn.isPositive(test_makeUnhappyNotes()) {
		t.Error("this should have been recoreded as negitive")
	}

}
