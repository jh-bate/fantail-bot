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
		s *session
		u Users
	}
)

func NewFollowUp(s *session) *FollowUp {
	sched := &FollowUp{
		s: s,
		c: cron.New(),
		u: Users{},
	}

	sched.Setup([]Task{GatherTask{}, RemindersTask{}, HelpMeTask{}, YouThereTask{}})

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
		users, err := fu.s.Storage.GetUsers()
		if err != nil {
			log.Println("Trying to run scheduled task ", err.Error())
			log.Println("Will bail ...")
			return
		}

		for i := range users {

			n, err := shed.s.Storage.GetLatest(users[i], 10)

			if err != nil {
				log.Println("Error getting latest ", err.Error())
				break
			}

			user := shed.u.GetUser(users[i])
			if user == nil {
				user = &User{}
			}
			user.id = users[i]
			if len(n) > 0 {
				user.recent = n
				user.lastChat = n.SortByDate()[0].AddedOn
			}
		}
		return
	}
}

func (this *GatherTask) spec() string {
	return "@daily"
}

func (this *RemindersTask) run(fu *FollowUp) func() {
	return func() {
		now := time.Now()
		for i := range fu.u {
			for r := range fu.u[i].GetReminders() {
				reminder := fu.u[i].GetReminders()[r]
				if reminder.RemindToday() {
					fu.s.User = fu.u[i].ToBotUser()
					shed.s.send(reminder.Text)
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
		for i := range shed.u {

			fu.s.User = fu.u[i].ToBotUser()

			last := fu.u[i].lastChat
			if now.YearDay()-last.YearDay() > 7 {
				shed.s.send(fmt.Sprintf("Long time no chat! Wanna %s or %s something?", chat_action, say_action))
			}
		}
	}
}

func (this *YouThereTask) spec() string {
	return "@daily"
}
