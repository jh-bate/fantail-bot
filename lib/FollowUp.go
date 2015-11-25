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
		log.Println("Running gather info ....")
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
				log.Println("adding notes for user")
				user.lastChat = n.MostRecent().AddedOn
				log.Println("set most recent as", user.lastChat)
			}

			fu.users = user.AddOrUpdate(fu.users)
		}
		return
	}
}

func (this *GatherTask) spec() string {
	return "0 0/5 * * *"
}

func (this *RemindersTask) run(fu *FollowUp) func() {
	return func() {
		log.Printf("Running reminders for %d users", len(fu.users))
		for i := range fu.users {
			log.Println("quick hi", fu.users[i].lastChat)
			fu.session.User = fu.users[i].ToBotUser()
			for r := range fu.users[i].GetReminders() {
				log.Printf("User has reminders ...")
				reminder := fu.users[i].GetReminders()[r]
				if reminder.RemindToday() {
					fu.session.User = fu.users[i].ToBotUser()
					fu.session.send(reminder.Text)
					//complete the reminder and save the update
					updated := reminder
					updated.CompletedOn = time.Now()
					fu.session.Storage.Update(string(fu.users[i].id), *reminder, *updated)
				}
			}

		}
	}
}

func (this *RemindersTask) spec() string {
	return "0 0/5 * * *"
}

func (this *HelpMeTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running `Help me` ....")
		return
	}
}

func (this *HelpMeTask) spec() string {
	return "0 0/5 * * *"
}

func (this *YouThereTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running `you there?` ....")
		now := time.Now()
		for i := range fu.users {

			fu.session.User = fu.users[i].ToBotUser()

			log.Println("User last chated", fu.users[i].lastChat)

			last := fu.users[i].lastChat
			if now.YearDay()-last.YearDay() > 3 {
				fu.session.send(fmt.Sprintf("Long time no chat! Wanna %s or %s something?", chat_action, say_action))
				//remove the users rescent history
				fu.users[i].recent = nil
			}
		}
	}
}

func (this *YouThereTask) spec() string {
	return "0 0/5 * * *"
}
