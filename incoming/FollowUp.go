package incoming

import (
	"fmt"
	"log"

	"github.com/jh-bate/fantail-bot/note"
	"github.com/jh-bate/fantail-bot/user"

	"github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/robfig/cron"
)

type (
	Task interface {
		run(f *FollowUp) func()
		// spec notes:
		//
		// *    *     *     *   *    *      command to be executed
		// -    -     -     -   -    -
		// |    |     |     |   |    |
		// |    |     |     |   |    +----- day of week (0 - 6) (Sunday=0)
		// |    |     |     |   +------- month (1 - 12)
		// |    |     |     +--------- day of month (1 - 31)
		// |    |     +----------- hour (0 - 23)
		// |     +------------- min (0 - 59)
		// +------------- sec (0 - 59)
		spec() string
	}

	Tasks []Task

	GatherTask    struct{}
	FollowupTask  struct{}
	CheckInTask   struct{}
	LearnFromTask struct{}

	FollowUp struct {
		*cron.Cron
		*Session
		user.Users
	}
)

func NewFollowUp(s *Session) *FollowUp {
	f := &FollowUp{
		Session: s,
		Cron:    cron.New(),
		Users:   user.Users{},
	}
	f.setup([]Task{&GatherTask{}, &FollowupTask{}, &CheckInTask{}, &LearnFromTask{}})
	return f
}

func (this *FollowUp) setup(t Tasks) {
	for i := range t {
		this.Cron.AddFunc(t[i].spec(), t[i].run(this))
	}
}

func (this *FollowUp) Start() {
	this.Cron.Start()
	return
}

func (this *FollowUp) Stop() {
	this.Cron.Stop()
	return
}

func (this *GatherTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running gather info ....")
		users, err := user.GetUsers()
		if err != nil {
			log.Println("Trying to run scheduled task ", err.Error())
			log.Println("Will bail ...")
			return
		}

		for i := range users {

			user := users.GetUser(users[i].Id)
			//if user == nil {
			//	user = &user.
			//}
			user = users[i]

			notes, err := note.GetNotes(string(users[i].Id))

			if err != nil {
				log.Println("Error getting latest ", err.Error())
				break
			}
			if len(notes) > 0 {
				user.Notes = notes.SortByDate()
			}

			fu.Users = user.AddOrUpdate(fu.Users)
		}
		return
	}
}

func (this *GatherTask) spec() string {
	//Every 5 mins
	return "0 0/5 * * *"
}

func (this *FollowupTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running `Help me` ....")
		for i := range fu.Users {

			help := fu.Users[i].NeedsHelp()

			if len(help) > 0 {
				fu.send(
					fu.Users[i].ToBotUser(),
					fmt.Sprintf(
						"Hey, so these are the things you wanted help with /n/n%s",
						help.ToString(),
					),
				)
			}
		}
		return
	}
}

func (this *FollowupTask) spec() string {
	//Every 10 mins
	return "0 0/10 * * *"
}

func (this *CheckInTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running `you there?` ....")
		for i := range fu.Users {

			keyboard := Keyboard{}
			keyboard = append(keyboard, []string{"/say all good thanks"}, []string{"/chat sounds like good idea"})

			fu.sendWithKeyboard(
				fu.Users[i].ToBotUser(),
				fmt.Sprintf("Long time no chat! Wanna %s or %s something?", chat_action, say_action),
				keyboard,
			)
		}
	}
}

func (this *CheckInTask) spec() string {
	//every 12 hours between 7am-8pm
	return "0 0 7-20/12 * * *"
}

func (this *LearnFromTask) run(fu *FollowUp) func() {

	const check_for_days = 3

	return func() {
		log.Println("Running `learning task` ....")
		for i := range fu.Users {

			pos := fu.Users[i].LearnAbout(check_for_days)
			keyboard := Keyboard{}

			if !pos {

				keyboard = append(keyboard, []string{"/say yeah things aren't going well"}, []string{"/say actually it is going well"})

				fu.sendWithKeyboard(
					fu.Users[i].ToBotUser(),
					"Hey, looks like things might not be going as well as you would like?",
					keyboard,
				)
				break
			}

			keyboard = append(keyboard, []string{"/say yeah I am doing well!"}, []string{"/say actually its not going well"})

			fu.sendWithKeyboard(
				fu.Users[i].ToBotUser(),
				"Hey, it looks like your doing well!!",
				keyboard,
			)

		}
	}
}

func (this *LearnFromTask) spec() string {
	//7am on MON,WED,FRI,SUN
	return "0 0 6 'MON,WED,FRI,SUN' * *"
}