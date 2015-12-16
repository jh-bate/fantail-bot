package sticker

import "github.com/jh-bate/fantail-bot/config"

type (
	Sticker struct {
		Ids     []string `json:"ids"`
		Meaning string   `json:"meaning"`
		SaveTag string   `json:"saveTag"`
	}

	Stickers []*Sticker
)

func Load() Stickers {
	var s Stickers
	config.Load(&s, "stickers.json")
	return s
}

func (this Stickers) Find(id string) *Sticker {
	for i := range this {
		for si := range this[i].Ids {
			if this[i].Ids[si] == id {
				return this[i]
			}
		}
	}
	return nil
}
