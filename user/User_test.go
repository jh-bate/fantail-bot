package user_test

import (
	"strings"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"
	. "github.com/jh-bate/fantail-bot/user"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("User", func() {

	var (
		myUser *User
	)

	const (
		userid = 999
	)

	BeforeEach(func() {
		myUser = New(userid)
	})

	Describe("When created", func() {
		It("should have the id set", func() {
			Expect(myUser.Id).To(Equal(userid))
		})

		It("should be able to convert to a BotUser", func() {
			botUser := myUser.ToBotUser()
			Expect(botUser.ID).To(Equal(myUser.Id))
			var botType telebot.User
			Expect(botUser).To(BeAssignableToTypeOf(botType))
		})
	})

	Describe("When created", func() {
		It("should have the id set", func() {
			Expect(myUser.Id).To(Equal(userid))
		})

		It("should be able to convert to a BotUser", func() {
			botUser := myUser.ToBotUser()
			Expect(botUser.ID).To(Equal(myUser.Id))
			var botType telebot.User
			Expect(botUser).To(BeAssignableToTypeOf(botType))
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
		u1_id = 667
		u2_id = 8868
		u3_id = 999
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

	Describe("When calling AddOrUpdate", func() {
		It("should add a user to the list if new", func() {
			Expect(len(myUsers)).To(Equal(3))
			other := New(22899991)
			myUsers = other.AddOrUpdate(myUsers)
			Expect(len(myUsers)).To(Equal(4))
		})
		It("should update an existing user if they are already in the list", func() {
			Expect(len(myUsers)).To(Equal(3))
			u3.Learnt = append(u3.Learnt, Learning{Date: time.Now(), Positive: true, Period: 5})

			myUsers = u3.AddOrUpdate(myUsers)
			Expect(len(myUsers)).To(Equal(3))

			updated := myUsers.GetUser(u3_id)
			Expect(len(updated.Learnt)).To(Equal(1))
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
