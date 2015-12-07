package lib

import (
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
)

type (
	Sticker struct {
		Ids     []string `json:"ids"`
		Meaning string   `json:"meaning"`
		SaveTag string   `json:"saveTag"`
	}

	Stickers []*Sticker
)

func LoadKnownStickers() Stickers {
	var s Stickers
	LoadConfig(&s, "stickers.json")
	return s
}

func (this Sticker) ToNote(msg telebot.Message, tags ...string) Note {

	tags = append(tags, this.Ids...)
	tags = append(tags, this.SaveTag)

	return Note{
		UserId: msg.Sender.ID,
		Added:  msg.Time(),
		Text:   this.Meaning,
		Tag:    strings.Join(tags, ","),
		Remind: time.Now().AddDate(0, 0, 7)}

}

func (this Stickers) FindSticker(id string) *Sticker {
	for i := range this {
		for si := range this[i].Ids {
			if this[i].Ids[si] == id {
				return this[i]
			}
		}
	}
	return nil
}
