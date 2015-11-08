package lib

type Action interface {
	getName() string
	getHint() string
}

const (
	say_action      = "/say"
	say_action_hint = "/say [what you want to say]"

	chat_action = "/chat"

	review_action      = "/review"
	review_action_hint = "/review <days>"

	remind_action      = "/remind"
	remind_action_hint = "/remind in <days> to <msg>"
)

func NewAction(in Incoming) Action {

	if in.getCmd() == say_action {
		return &SayAction{}
	} else if in.getCmd() == remind_action {
		return &RemindAction{}
	} else if in.getCmd() == review_action {
		return &ReviewAction{}
	} else if in.getCmd() == chat_action {
		return &ChatAction{}
	}
	return nil
}

type SayAction struct{}

func (a SayAction) getName() string {
	return say_action
}
func (a SayAction) getHint() string {
	return say_action_hint
}

type ChatAction struct{}

func (a ChatAction) getName() string {
	return chat_action
}
func (a ChatAction) getHint() string {
	return ""
}

type ReviewAction struct{}

func (a ReviewAction) getName() string {
	return review_action
}
func (a ReviewAction) getHint() string {
	return review_action_hint
}

func (a ReviewAction) getDays() int {
	return 0
}

type RemindAction struct{}

func (a RemindAction) getName() string {
	return remind_action
}
func (a RemindAction) getHint() string {
	return remind_action_hint
}

func (a RemindAction) getDays() int {
	return 0
}

func (a RemindAction) getMessage() string {
	return ""
}
