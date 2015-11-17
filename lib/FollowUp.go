package lib

import (
	"fmt"
	"log"
	"time"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/robfig/cron"
)

type (
	Task interface {
		run(f *FollowUp) func()
		spec() string
	}

	Tasks []Task

	GatherTask    struct{}
	RemindersTask struct{}
	HelpMeTask    struct{}
	YouThereTask  struct{}

	FollowUp struct {
		c *cron.Cron
		*session
		users Users
	}
)

func NewFollowUp(s *session) *FollowUp {
	sched := &FollowUp{
		session: s,
		c:       cron.New(),
		users:   Users{},
	}
	sched.setup([]Task{&GatherTask{}, &RemindersTask{}, &HelpMeTask{}, &YouThereTask{}})

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
		//when did we last chat?
		users, err := fu.session.Storage.GetUsers()
		if err != nil {
			log.Println("Trying to run scheduled task ", err.Error())
			log.Println("Will bail ...")
			return
		}

		for i := range users {

			n, err := fu.session.Storage.GetLatest(users[i], 10)

			if err != nil {
				log.Println("Error getting latest ", err.Error())
				break
			}

			user := fu.users.GetUser(users[i])
			if user == nil {
				user = &User{}
			}
			user.id = users[i]
			if len(n) > 0 {
				user.recent = n
				user.lastChat = n.SortByDate()[0].AddedOn
			}
			fu.users.AddUser(user)
		}
		return
	}
}

func (this *GatherTask) spec() string {
	return "@daily"
}

func (this *RemindersTask) run(fu *FollowUp) func() {
	return func() {
		for i := range fu.users {
			for r := range fu.users[i].GetReminders() {
				reminder := fu.users[i].GetReminders()[r]
				if reminder.RemindToday() {
					fu.session.User = fu.users[i].ToBotUser()
					fu.session.send(reminder.Text)
				}
			}

		}
	}
}

func (this *RemindersTask) spec() string {
	return "0 30"
}

func (this *HelpMeTask) run(fu *FollowUp) func() {
	return func() {
		return
	}
}

func (this *HelpMeTask) spec() string {
	return "@daily"
}

func (this *YouThereTask) run(fu *FollowUp) func() {
	return func() {
		now := time.Now()
		for i := range fu.users {

			fu.session.User = fu.users[i].ToBotUser()

			last := fu.users[i].lastChat
			if now.YearDay()-last.YearDay() > 7 {
				fu.session.send(fmt.Sprintf("Long time no chat! Wanna %s or %s something?", chat_action, say_action))
			}
		}
	}
}

func (this *YouThereTask) spec() string {
	return "@daily"
}
