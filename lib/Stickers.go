package lib

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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

	dirPath, err := os.Getwd()
	if err != nil {
		log.Panic("could not working dir ", err.Error())
	}

	var s Stickers
	file, err := os.Open(fmt.Sprintf("%s/config/stickers.json", dirPath))
	if err != nil {
		log.Panic("could not load stickers config ", err.Error())
	}
	err = json.NewDecoder(file).Decode(&s)
	if err != nil {
		log.Panic("could not load stickers config ", err.Error())
	}
	return s
}

func (this Sticker) ToNote(msg telebot.Message, tags ...string) Note {

	tags = append(tags, this.Ids...)
	tags = append(tags, this.SaveTag)

	return Note{
		WhoId:      msg.Sender.ID,
		AddedOn:    msg.Time(),
		Text:       this.Meaning,
		Tag:        strings.Join(tags, ","),
		RemindNext: time.Now().AddDate(0, 0, 7)}

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
