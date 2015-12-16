package question_test

import (
	"encoding/json"
	"log"

	. "github.com/jh-bate/fantail-bot/question"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Question", func() {

	var (
		q struct {
			Questions `json:"QandA"`
		}

		qJson = []byte(`
	{
	  "QandA": [
	  	{
	      "relatesTo": null,
	      "context": [
	        "Q1"
	      ],
	      "question": "Q1?",
	      "answers": [
	        "Q1 A1",
	        "Q1 A2"
	      ]
	    },
	    {
	      "relatesTo": {
	        "answers": [
	          "Q1 A1",
	        "Q1 A2"
	        ],
	        "save": true,
	        "saveTag": "Q1"
	      },
	      "context": [
	        "Q2"
	      ],
	      "question": "Q2?",
	      "answers": [
	        "Q2 A1",
	        "Q2 A2"
	      ]
	    },
	    {
	      "relatesTo": {
	        "answers": [
	          "Q2 A1",
	          "Q2 A2"
	        ],
	        "save": true,
	        "saveTag": "Q2"
	      },
	      "context": [
	        "Q3"
	      ],
	      "question": "Q3?",
	      "answers": [
	        "Q3 A1",
	        "Q3 A2"
	      ]
	    }
	  ]
	}
	`)
	)

	BeforeEach(func() {

		err := json.Unmarshal(qJson, &q)
		if err != nil {
			log.Panic("could not decode QandA ", err.Error())
		}

	})

	Describe("When Loaded", func() {
		It("first quesion should be set", func() {
			Expect(q.Questions.First()).To(Not(BeNil()))
		})
		It("first quesion should be a question", func() {
			var qType *Question
			Expect(q.Questions.First()).To(BeAssignableToTypeOf(qType))
		})
		It("can find the next question based on the given answer", func() {
			n, save := q.Questions.Next(q.Questions[0].PossibleAnswers[0])
			Expect(save).To(BeTrue())
			Expect(n).To(Equal(q.Questions[1]))
		})
		It("can find the next question based on a range of answers", func() {
			question, _ := q.Questions.NextFrom(q.Questions.First().PossibleAnswers...)
			Expect(question).To(Equal(q.Questions[1]))
		})
		It("can make a keyboad from the given answer", func() {

			q1 := q.Questions.First()
			kb := q1.MakeKeyboard()

			for i := range q1.PossibleAnswers {
				if q1.PossibleAnswers[i] != kb[i][0] {
					Expect(q1.PossibleAnswers[i]).To(Equal(kb[i][0]))
				}
			}
		})

	})
	Describe("When Empty", func() {
		It("first question will not be set", func() {
			var noQ Questions
			Expect(noQ.First()).To(BeNil())
		})
	})

})
