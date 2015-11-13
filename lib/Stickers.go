package lib

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
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

	absPath, _ := filepath.Abs("config/stickers.json")
	log.Println("sticker path ", absPath)

	file, err := os.Open(absPath)
	if err != nil {
		log.Println("could not load stickers config ", err.Error())
		absPath, _ = filepath.Abs("lib/config/stickers.json")
		log.Println("sticker path ", absPath)

		file, err = os.Open(absPath)
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
