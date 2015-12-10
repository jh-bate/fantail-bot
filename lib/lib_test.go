package lib_test

import (
	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
	. "github.com/jh-bate/fantail-bot/lib"

	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("Lib", func() {

	Describe("Note", func() {
		var (
			note *Note
		)

		BeforeEach(func() {

			note = NewNote(telebot.Message{
				Text:   "/say hi",
				Sender: telebot.User{FirstName: "my user", ID: 12345},
			}, "testing")
		})

		Context("New note", func() {
			It("should be current", func() {
				Expect(note.IsCurrent()).To(Equal(true))
			})
		})
	})
})
