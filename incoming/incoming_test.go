package incoming //note also testing unexported methods

import (
	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
	"github.com/jh-bate/fantail-bot/question"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Payload", func() {

	var (
		standardMsg       telebot.Message
		stdPayload        *Payload
		submissionMsg     telebot.Message
		submissionPayload *Payload
		stickerMsg        telebot.Message
		stickerPayload    *Payload
	)

	const (
		msg_username       = "My tester"
		standard_msg_text  = "my message"
		submisson_msg_text = "/action and stuff"
		sticker_id         = "BQADAwADCgADt6a9BtopSv1uQpPwAg"
		msg_userid         = 9992476
	)

	BeforeEach(func() {

		standardMsg = telebot.Message{
			Text:   standard_msg_text,
			Sender: telebot.User{FirstName: msg_username, ID: msg_userid},
		}

		stdPayload = New(standardMsg)

		submissionMsg = telebot.Message{
			Text:   submisson_msg_text,
			Sender: telebot.User{FirstName: msg_username, ID: msg_userid},
		}

		submissionPayload = New(submissionMsg)

		file, _ := telebot.NewFile("./payload.go")
		file.FileID = sticker_id

		stickerMsg = telebot.Message{
			Text:    "",
			Sender:  telebot.User{FirstName: msg_username, ID: msg_userid},
			Sticker: telebot.Sticker{File: file},
		}

		stickerPayload = New(stickerMsg)

	})

	Describe("When created", func() {

		It("should have a date the same as original message", func() {
			Expect(stdPayload.Date).To(Equal(standardMsg.Time()))
		})
		It("should have a text the same as original message", func() {
			Expect(stdPayload.Text).To(Equal(standardMsg.Text))
		})
		It("should have an action is it is set from original message", func() {
			Expect(stdPayload.Action).To(Equal(""))
		})
		It("should have not have a submission", func() {
			Expect(stdPayload.HasSubmisson).To(BeFalse())
		})
		It("should have not have a sticker", func() {
			Expect(stdPayload.Sticker.Exists).To(BeFalse())
		})
		It("should not have an action", func() {
			Expect(stdPayload.HasAction()).To(BeFalse())
		})
	})

	Describe("When created as submission", func() {
		It("should have a date the same as original message", func() {
			Expect(submissionPayload.Date).To(Equal(submissionMsg.Time()))
		})
		It("should have a text the same as original message", func() {
			Expect(submissionPayload.Text).To(Equal(submissionMsg.Text))
		})
		It("should have an action", func() {
			Expect(submissionPayload.Action).To(Equal("/action"))
		})
		It("should have have a submission", func() {
			Expect(submissionPayload.HasSubmisson).To(BeTrue())
		})
		It("should have not have a sticker", func() {
			Expect(submissionPayload.Sticker.Exists).To(BeFalse())
		})
		It("should have an action", func() {
			Expect(submissionPayload.HasAction()).To(BeTrue())
		})
	})

	Describe("When created as a sticker", func() {
		It("should have a date the same as original message", func() {
			Expect(stickerPayload.Date).To(Equal(stickerMsg.Time()))
		})
		It("should have a text the same as original message", func() {
			Expect(stickerPayload.Text).To(Equal(stickerMsg.Text))
		})
		It("should have an action", func() {
			Expect(stickerPayload.Action).To(Equal(""))
		})
		It("should have have a submission", func() {
			Expect(stickerPayload.HasSubmisson).To(BeFalse())
		})
		It("should have a sticker", func() {
			Expect(stickerPayload.Sticker.Exists).To(BeTrue())
		})
		It("should have a sticker", func() {
			Expect(stickerPayload.Sticker.Id).To(Equal(sticker_id))
		})
		It("should not have an action", func() {
			Expect(stickerPayload.HasAction()).To(BeFalse())
		})
	})

})

var _ = Describe("Actions", func() {

	var (
		action Action
	)

	BeforeEach(func() {

		session := NewSession(&telebot.Bot{})

		submissionMsg := telebot.Message{
			Text:   "/say hi!",
			Sender: telebot.User{FirstName: "testing", ID: 999},
		}

		payload := New(submissionMsg)

		action = NewAction(payload, "", session)

	})

	Describe("When action created", func() {

		It("should have name be what was sent in the message", func() {
			Expect(action.getName()).To(Equal("/say"))
		})
		It("should get hint for action", func() {
			Expect(action.getHint()).To(Equal("/say <message>"))
		})
		It("should return question", func() {
			var question *question.Question
			Expect(action.nextQuestion()).To(BeAssignableToTypeOf(question))
		})
		It("should return questions", func() {
			var questions question.Questions
			Expect(action.getQuestions()).To(BeAssignableToTypeOf(questions))
		})
	})

})
