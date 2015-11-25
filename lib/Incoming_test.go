package lib

import (
	"strings"
	"testing"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

func newMsg(txt string) telebot.Message {

	return telebot.Message{
		Text:   txt,
		Sender: telebot.User{FirstName: "my user", ID: 12345},
	}
}

func newStickerMsg() telebot.Message {

	file, _ := telebot.NewFile("./config/fantail.json")
	file.FileID = "BQADAwADCgADt6a9BtopSv1uQpPwAg"

	sticker := telebot.Sticker{File: file}

	return telebot.Message{
		Text:    "",
		Sender:  telebot.User{FirstName: "my user", ID: 12345},
		Sticker: sticker,
	}
}

func TestIncoming_isCmd(t *testing.T) {

	inCmd := newIncoming(newMsg("/stuff do to"))

	if inCmd.isCmd() == false {
		t.Fail()
	}

	inString := newIncoming(newMsg("stuff do to"))

	if inString.isCmd() == true {
		t.Fail()
	}

	inEmpty := newIncoming(newMsg(""))

	if inEmpty.isCmd() == true {
		t.Fail()
	}

	inSticker := newIncoming(newStickerMsg())

	if inSticker.isCmd() == true {
		t.Fail()
	}

}

func TestIncoming_hasSubmisson(t *testing.T) {

	inCmd := newIncoming(newMsg("/stuff do to"))

	if inCmd.hasSubmisson() == false {
		t.Error("command with text should be a submisson")
	}

	inString := newIncoming(newMsg("stuff do to"))

	if inString.hasSubmisson() == true {
		t.Error("plain txt message should NOT a submisson")
	}

	inSticker := newIncoming(newStickerMsg())

	if inSticker.hasSubmisson() == true {
		t.Error("sticker message should NOT be a submisson")
	}

}

func TestIncoming_isSticker(t *testing.T) {

	inCmd := newIncoming(newMsg("/stuff do to"))

	if inCmd.isSticker() == true {
		t.Error("command with text is not a sticker")
	}

	inString := newIncoming(newMsg("stuff do to"))

	if inString.isSticker() == true {
		t.Error("text is not a sticker")
	}

	inSticker := newIncoming(newStickerMsg())

	if inSticker.isSticker() == false {
		t.Error("should be a sticker")
	}

}

func TestIncoming_sender(t *testing.T) {

	inCmd := newIncoming(newMsg("/stuff do to"))

	if inCmd.sender() != inCmd.msg.Sender {
		t.Error("should be the same as the telebot.Message.Sender")
	}

}

func TestIncoming_getNote(t *testing.T) {

	inCmd := newIncoming(newMsg("/stuff to do"))

	n := inCmd.getNote()

	if n.IsEmpty() {
		t.Error("should have created a note")
	}

	if n.Text != "to do" {
		t.Error("should have created a note")
	}

	if strings.Contains(n.Tag, "/stuff") == false {
		t.Error("should have created a note with /stuff tag")
	}

	inRemind := newIncoming(newMsg("/remind 3 to do some tests"))

	r := inRemind.getNote()

	if r.IsEmpty() {
		t.Error("should have created a note")
	}

	if r.Text != "to do some tests" {
		t.Error("should have created a note")
	}

	if strings.Contains(r.Tag, "/remind") == false {
		t.Error("should have created a note with /remind tag")
	}

}

func TestIncoming_getCmd(t *testing.T) {

	inCmd := newIncoming(newMsg("/stuff do to"))

	if inCmd.getCmd() != "/stuff" {
		t.Error("should be the command /stuff")
	}

	inNoCmd := newIncoming(newMsg("stuff do to"))

	if inNoCmd.getCmd() != "" {
		t.Error("should be no command")
	}

	inSticker := newIncoming(newStickerMsg())

	if inSticker.getCmd() != "" {
		t.Error("should be no command for a sticker")
	}

}

func TestIncoming_getCmd_returnsLowerCase(t *testing.T) {

	inCmd := newIncoming(newMsg("/Stuff do to"))

	if inCmd.getCmd() != "/stuff" {
		t.Error("should be the command /stuff")
	}

}

func TestIncoming_cmdMatches(t *testing.T) {

	inCmd := newIncoming(newMsg("/stuff do to"))

	if inCmd.cmdMatches("/no", "/match") {
		t.Error("should be no command match")
	}

	if inCmd.cmdMatches("/no", "/stuff") == false {
		t.Error("should match /stuff")
	}

}

func TestIncoming_cmdMatches_ingoresCase(t *testing.T) {

	inCmd := newIncoming(newMsg("/Stuff do to"))

	if !inCmd.cmdMatches("/no", "/stuff") {
		t.Error("should be a match on stuff")
	}

}

func TestIncoming_submissonMatches(t *testing.T) {

	inSub := newIncoming(newMsg("/stuff do to"))

	if inSub.submissonMatches("/no", "/match") {
		t.Error("should be no command match")
	}

	if inSub.submissonMatches("/no", "/stuff") == false {
		t.Error("should match /stuff")
	}

	inCmd := newIncoming(newMsg("/stuff"))

	if inCmd.submissonMatches("/stuff") {
		t.Error("should be no match as there isn't a submisson")
	}

}
