package note_test

import (
	"time"

	. "github.com/jh-bate/fantail-bot/note"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Note", func() {

	var (
		myNote *Note
	)

	const (
		test_note_text   = "testing testing 123"
		test_note_userid = 999
	)

	BeforeEach(func() {
		//note controller testing
		myNote = New(test_note_text, test_note_userid, time.Now(), SAID_TAG)
	})

	Describe("When created", func() {
		It("should be current", func() {
			Expect(myNote.IsCurrent()).To(Equal(true))
		})
		It("should not be empty", func() {
			Expect(myNote.IsEmpty()).To(Equal(false))
		})
		It("should have the userid set", func() {
			Expect(myNote.UserId).To(Equal(test_note_userid))
		})
		It("should have the text set", func() {
			Expect(myNote.Text).To(Equal(test_note_text))
		})
		It("should have the tag set", func() {
			Expect(myNote.Tag).To(Equal(SAID_TAG + ","))
		})
	})

	Describe("When updated", func() {

		It("should set the updated date", func() {
			myNote.Update()
			Expect(myNote.Updated.IsZero()).To(Equal(false))
		})
		It("should still be current", func() {
			myNote.Update()
			Expect(myNote.IsCurrent()).To(Equal(true))
		})
	})

	Describe("When completed", func() {

		It("should set the updated date", func() {
			myNote.Complete()
			Expect(myNote.Deleted.IsZero()).To(Equal(false))
		})
		It("should no longer be current", func() {
			myNote.Complete()
			Expect(myNote.IsCurrent()).To(Equal(false))
		})
	})

	/*Describe("When saved", func() {
		It("should not return an error", func() {
			Expect(myNote.Save()).To(BeNil())
		})
	})*/
	Describe("When deleted", func() {
		It("should not return an error", func() {
			Expect(myNote.Delete()).To(BeNil())
		})
	})

})

var _ = Describe("Notes", func() {

	var (
		myNotes Notes
		n1      *Note
		n2      *Note
		n3      *Note
	)

	const (
		test_note_text   = "testing testing 123"
		test_note_text_2 = "testing testing 456"
		test_note_text_3 = "testing testing 8910"
		test_note_userid = 999
	)

	BeforeEach(func() {
		n1 = New(test_note_text, test_note_userid, time.Now().AddDate(0, 0, -5), SAID_TAG)
		n2 = New(test_note_text_2, test_note_userid, time.Now(), HELP_TAG)
		n3 = New(test_note_text_3, test_note_userid, time.Now().AddDate(0, 0, 5), CHAT_TAG)
		myNotes = Notes{n1, n2, n3}
	})

	Describe("When sorted", func() {
		It("should be sorted by date", func() {
			sorted := myNotes.OldestFirst()
			Expect(sorted[0]).To(Equal(n1))
			Expect(sorted[1]).To(Equal(n2))
			Expect(sorted[2]).To(Equal(n3))
		})
	})

	Describe("When filtered by text", func() {
		It("should return only those that match", func() {
			notes := myNotes.FilterOnTxt(test_note_text_2)
			Expect(len(notes)).To(Equal(1))
			Expect(notes[0]).To(Equal(n2))
		})
	})

	Describe("When getting words", func() {
		It("should return all words in text and tags", func() {
			Expect(len(myNotes.GetWords())).To(Equal(12))
		})
	})

	Describe("When getting most recent", func() {
		It("should return only the newest note", func() {
			Expect(myNotes.MostRecent()).To(Equal(n3))
		})
		It("should an empty note when there are none", func() {
			empty := Notes{}
			Expect(empty.MostRecent().IsEmpty()).To(Equal(true))
		})
	})

	Describe("When getting newer than", func() {
		It("should return notes newer than x days ago", func() {
			Expect(len(myNotes.NewerThan(4))).To(Equal(2))
		})
	})

})
