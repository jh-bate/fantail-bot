package bot //note also testing unexported methods

import (
	"time"

	"github.com/jh-bate/fantail-bot/question"

	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("Payload", func() {

	var (
		stdPayload        *Payload
		stdTime           time.Time
		submissionPayload *Payload
		submissionTime    time.Time
	)

	const (
		standard_msg_text  = "my message"
		submisson_msg_text = "/action and stuff"
		msg_userid         = "9992476"
		msg_user_name      = "Me"
	)

	BeforeEach(func() {

		stdTime = time.Now().Add(-2)
		stdPayload = New(msg_userid, msg_user_name, standard_msg_text, stdTime)

		submissionTime = time.Now().Add(-5)
		submissionPayload = New(msg_userid, msg_user_name, submisson_msg_text, submissionTime)

	})

	Describe("When created", func() {

		It("should have a date the same as original message", func() {
			Expect(stdPayload.Date).To(Equal(stdTime))
		})
		It("should have a text the same as original message", func() {
			Expect(stdPayload.Text).To(Equal(standard_msg_text))
		})
		It("should have an action is it is set from original message", func() {
			Expect(stdPayload.Action).To(Equal(""))
		})
		It("should have not have a submission", func() {
			Expect(stdPayload.HasSubmisson).To(BeFalse())
		})
		It("should not have an action", func() {
			Expect(stdPayload.HasAction()).To(BeFalse())
		})
	})

	Describe("When created as submission", func() {
		It("should have a date the same as original message", func() {
			Expect(submissionPayload.Date).To(Equal(submissionTime))
		})
		It("should have a text the same as original message", func() {
			Expect(submissionPayload.Text).To(Equal(submisson_msg_text))
		})
		It("should have an action", func() {
			Expect(submissionPayload.Action).To(Equal("/action"))
		})
		It("should have have a submission", func() {
			Expect(submissionPayload.HasSubmisson).To(BeTrue())
		})
		It("should have an action", func() {
			Expect(submissionPayload.HasAction()).To(BeTrue())
		})
	})

})

var _ = Describe("Action", func() {

	var (
		action  Action
		payload *Payload
	)

	BeforeEach(func() {

		payload = New("999", "You", "/say hi!", time.Now())

		action = NewAction(
			payload,
			"",
			NewSession(nil),
		)

	})

	Describe("When created", func() {

		It("should have name be what was sent in the message", func() {
			Expect(action.Name()).To(Equal("/say"))
		})
		It("should return the payload it was initialised with", func() {
			Expect(action.Payload()).To(Equal(payload))
		})
	})

	Describe("When passed to next question", func() {
		It("should return a question for that action", func() {
			Expect(nextQuestion(action)).To(Not(BeNil()))
			var q *question.Question
			Expect(nextQuestion(action)).To(BeAssignableToTypeOf(q))
		})
	})

})
