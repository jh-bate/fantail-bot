package lib

import "testing"

func TestSticker_LoadKnownStickers(t *testing.T) {

	stickers := LoadKnownStickers()

	if len(stickers) == 0 {
		t.Error("Stickers should have been loaded")
	}
}

func TestSticker_FindSticker(t *testing.T) {

	const stickerId = "BQADAwADDAADt6a9BkUSLrxnvwHfAg"

	stickers := LoadKnownStickers()

	sLow := stickers.FindSticker(stickerId)

	if sLow == nil {
		t.Error("Sticker should have been loaded")
	}

	match := false
	for i := range sLow.Ids {
		if sLow.Ids[i] == stickerId {
			match = true
		}
	}

	if match == false {
		t.Errorf("Sticker id %s should have been in %v", stickerId, sLow.Ids)
	}
}
