package lib

import (
	"encoding/json"
	"log"
	"testing"
)

func loadQuestions() Questions {

	var q struct {
		Questions `json:"QandA"`
	}

	qJson := []byte(`
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

	err := json.Unmarshal(qJson, &q)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}
	return q.Questions
}

func TestQuestions_First(t *testing.T) {

	q := loadQuestions()
	if q.First() == nil {
		t.Error("should return a question")
	}

	var nQ Questions
	if nQ.First() != nil {
		t.Error("should NOT return a question")
	}
}

func TestQuestions_next(t *testing.T) {

	q := loadQuestions()

	n, _ := q.next(q[0].PossibleAnswers[0])

	if n != q[1] {
		t.Errorf("expected %v got %v", q[1], n)
	}

}
