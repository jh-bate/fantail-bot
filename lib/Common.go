package lib

import (
	"strings"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	Process interface {
		Run(input <-chan telebot.Message)
	}
)

func hasSubmisson(txt string, cmds ...string) bool {

	if isCmd(txt, cmds...) {
		words := strings.Fields(txt)
		if len(words) > 1 {
			return true
		}
	}
	return false
}

func isCmd(txt string, cmds ...string) bool {

	for i := range cmds {
		if strings.Contains(txt, cmds[i]) {
			return true
		}
	}
	return false
}
