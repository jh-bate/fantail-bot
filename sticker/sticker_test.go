package sticker_test

import (
	. "github.com/jh-bate/fantail-bot/sticker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sticker", func() {

	var (
		myStickers Stickers
	)

	BeforeEach(func() {
		myStickers = Load()
	})

	Describe("When loaded", func() {
		It("should be stickers", func() {
			Expect(len(myStickers) > 0).To(Equal(true))
		})
	})

	Describe("When finding", func() {

		const stickerId = "BQADAwADDAADt6a9BkUSLrxnvwHfAg"

		It("should should be found", func() {
			theSticker := myStickers.Find(stickerId)
			Expect(theSticker).To(Not(BeNil()))
		})

		It("should should contain the id", func() {
			theSticker := myStickers.Find(stickerId)

			match := false
			for i := range theSticker.Ids {
				if theSticker.Ids[i] == stickerId {
					match = true
				}
			}

			Expect(match).To(BeTrue())
		})
	})

})
