package incoming

import (
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

	file, _ := telebot.NewFile("./config/default.json")
	file.FileID = "BQADAwADCgADt6a9BtopSv1uQpPwAg"

	sticker := telebot.Sticker{File: file}

	return telebot.Message{
		Text:    "",
		Sender:  telebot.User{FirstName: "my user", ID: 12345},
		Sticker: sticker,
	}
}

func TestMessage_Action(t *testing.T) {

	in := New(newMsg("/stuff do to"))

	if in.Action != "/stuff" {
		t.Fail()
	}

	in2 := newIncoming(newMsg("stuff do to"))

	if in2.Action != "" {
		t.Fail()
	}

	in3 := newIncoming(newMsg(""))

	if in3.Action != "" {
		t.Fail()
	}

	inSticker := newIncoming(newStickerMsg())

	if inSticker.Action != "" {
		t.Fail()
	}

	inCmd := New(newMsg("/STUFF do to"))

	if inCmd.Action != "/stuff" {
		t.Error("should be lower case /stuff")
	}
}

func TestMessage_HasSubmisson(t *testing.T) {

	inCmd := New(newMsg("/stuff do to"))

	if inCmd.HasSubmisson == false {
		t.Error("command with text should be a submisson")
	}

	inString := New(newMsg("stuff do to"))

	if inString.HasSubmisson == true {
		t.Error("plain txt message should NOT a submisson")
	}

	inSticker := New(newStickerMsg())

	if inSticker.HasSubmisson == true {
		t.Error("sticker message should NOT be a submisson")
	}

}

func TestMessage_HasSticker(t *testing.T) {

	inCmd := New(newMsg("/stuff do to"))

	if inCmd.HasSticker == true {
		t.Error("command with text is not a sticker")
	}

	inString := New(newMsg("stuff do to"))

	if inString.HasSticker == true {
		t.Error("text is not a sticker")
	}

	inSticker := New(newStickerMsg())

	if inSticker.HasSticker == false {
		t.Error("should be a sticker")
	}

}

func TestMessage_Sent(t *testing.T) {

	inCmd := New(newMsg("/stuff do to"))

	if inCmd.Sent != inCmd.msg.Sender {
		t.Error("should be the same as the telebot.Message.Sender")
	}
}
