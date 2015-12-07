package lib

import (
	"fmt"
	"log"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/robfig/cron"
)

type (
	Task interface {
		run(f *FollowUp) func()
		spec() string
	}

	Tasks []Task

	GatherTask    struct{}
	FollowupTask  struct{}
	CheckInTask   struct{}
	LearnFromTask struct{}

	FollowUp struct {
		c *cron.Cron
		*session
		users Users
	}
)

// Scheduling notes:
//
// *    *     *     *   *    *        command to be executed
// -    -     -     -   -    -
// |    |     |     |   |    |
// |    |     |     |   |    +----- day of week (0 - 6) (Sunday=0)
// |    |     |     |   +------- month (1 - 12)
// |    |     |     +--------- day of month (1 - 31)
// |    |     +----------- hour (0 - 23)
// |     +------------- min (0 - 59)
// +------------- sec (0 - 59)
func NewFollowUp(s *session) *FollowUp {
	sched := &FollowUp{
		session: s,
		c:       cron.New(),
		users:   Users{},
	}
	sched.setup([]Task{&GatherTask{}, &FollowupTask{}, &CheckInTask{}, &LearnFromTask{}})

	return sched
}

func (this *FollowUp) setup(t Tasks) {
	for i := range t {
		this.c.AddFunc(t[i].spec(), t[i].run(this))
	}
}

func (this *FollowUp) Start() {
	this.c.Start()
	return
}

func (this *FollowUp) Stop() {
	this.c.Stop()
	return
}

func (this *GatherTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running gather info ....")
		users, err := fu.session.Storage.GetUsers()
		if err != nil {
			log.Println("Trying to run scheduled task ", err.Error())
			log.Println("Will bail ...")
			return
		}

		for i := range users {

			user := fu.users.GetUser(users[i])
			if user == nil {
				user = &User{}
			}
			user.id = users[i]

			n, err := fu.session.Storage.Get(string(users[i]))

			if err != nil {
				log.Println("Error getting latest ", err.Error())
				break
			}
			if len(n) > 0 {
				user.notes = n.SortByDate()
			}

			fu.users = user.AddOrUpdate(fu.users)
		}
		return
	}
}

func (this *GatherTask) spec() string {
	//every 5 mins
	return "0 0/5 * * *"
}

func (this *FollowupTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running `Help me` ....")
		for i := range fu.users {

			help := fu.users[i].FollowUpAbout()

			if len(help) > 0 {
				fu.session.User = fu.users[i].ToBotUser()
				helpTxt := help.ToString()
				fu.session.send(fmt.Sprintf("Hey, so these are the things you wanted help with /n/n%s", helpTxt))
			}
		}
		return
	}
}

func (this *FollowupTask) spec() string {
	//every 10 mins
	return "0 0/10 * * *"
}

func (this *CheckInTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running `you there?` ....")
		for i := range fu.users {
			fu.session.User = fu.users[i].ToBotUser()

			keyboard := Keyboard{}
			keyboard = append(keyboard, []string{"/say all good thanks"}, []string{"/chat sounds like good idea"})

			fu.session.sendWithKeyboard(
				fmt.Sprintf("Long time no chat! Wanna %s or %s something?", chat_action, say_action),
				keyboard,
			)
		}
	}
}

func (this *CheckInTask) spec() string {
	//every day at 7am
	return "0 0 6 * * *"
}

func (this *LearnFromTask) run(fu *FollowUp) func() {

	const check_for_days = 3

	return func() {
		log.Println("Running `learning task` ....")
		for i := range fu.users {
			fu.session.User = fu.users[i].ToBotUser()

			pos := fu.users[i].IsPostive(check_for_days)
			keyboard := Keyboard{}

			if !pos {

				keyboard = append(keyboard, []string{"/say yeah things aren't going well"}, []string{"/say actually it is going well"})

				fu.session.sendWithKeyboard(
					"Hey, looks like things might not be going as well as you would like?",
					keyboard,
				)
				break
			}

			keyboard = append(keyboard, []string{"/say yeah I am doing well!"}, []string{"/say actually its not going well"})

			fu.session.sendWithKeyboard(
				"Hey, it looks like your doing well!!",
				keyboard,
			)

		}
	}
}

func (this *LearnFromTask) spec() string {
	//7pm on MON,WED,FRI,SUN
	return "0 0 19 'MON,WED,FRI,SUN' *"
}
