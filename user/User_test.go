package user_test

import (
	"strings"

	. "github.com/jh-bate/fantail-bot/user"

	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/onsi/gomega"
)

var _ = Describe("User", func() {

	var (
		myUser *User
	)

	const (
		userid = "999"
	)

	BeforeEach(func() {
		myUser = New(userid)
	})

	Describe("When created", func() {
		It("should have the id set", func() {
			Expect(myUser.Id).To(Equal(userid))
		})

	})

})

var _ = Describe("Users", func() {

	var (
		myUsers Users
		u1      *User
		u2      *User
		u3      *User
	)

	const (
		u1_id = "667"
		u2_id = "8868"
		u3_id = "999"
	)

	BeforeEach(func() {
		u1 = New(u1_id)
		u2 = New(u2_id)
		u3 = New(u3_id)
		myUsers = Users{u1, u2, u3}
	})

	Describe("When calling GetUser", func() {
		It("should find the user asked for", func() {
			user := myUsers.GetUser(u2_id)
			Expect(user).To(Equal(u2))
		})
	})
})

var _ = Describe("Classify", func() {

	var (
		classify *Classify
	)

	BeforeEach(func() {
		classify = NewClassification()
	})

	Describe("When happy words", func() {

		var happy []string
		happy = append(happy, strings.Fields("/say really happy")...)
		happy = append(happy, strings.Fields("/say hello all good")...)
		happy = append(happy, strings.Fields("things very good, went really well today")...)
		happy = append(happy, strings.Fields("/say bad day today, too many lows")...)

		It("should be positive", func() {
			Expect(classify.ArePositive(happy)).To(BeTrue())
		})
	})

	Describe("When unhappy words", func() {

		var unhappy []string
		unhappy = append(unhappy, strings.Fields("/say help")...)
		unhappy = append(unhappy, strings.Fields("/say more highs, sick of it!")...)
		unhappy = append(unhappy, strings.Fields("things went really well today")...)
		unhappy = append(unhappy, strings.Fields("/say all good help")...)
		unhappy = append(unhappy, strings.Fields("/say low again!!")...)

		It("should be negative", func() {
			Expect(classify.ArePositive(unhappy)).To(BeFalse())
		})
	})

	Describe("When neutral words", func() {

		var neutral []string
		neutral = append(neutral, strings.Fields("/say low happy")...)
		neutral = append(neutral, strings.Fields("/say bad happy")...)
		neutral = append(neutral, strings.Fields("all going well bad")...)
		neutral = append(neutral, strings.Fields("/say help good")...)
		neutral = append(neutral, strings.Fields("/say low great")...)

		It("should be positive", func() {
			Expect(classify.ArePositive(neutral)).To(BeTrue())
		})
	})

})
