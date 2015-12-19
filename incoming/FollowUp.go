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

	FollowupTask  struct{}
	CheckInTask   struct{}
	LearnFromTask struct{}

	FollowUp struct {
		*cron.Cron
		*Session
	}
)

func LoadUsersAndNotes() user.Users {

	log.Println("getting users ...")
	users, err := user.GetUsers()
	if err != nil {
		log.Println("Trying to run scheduled task ", err.Error())
		log.Println("bailing ...")
		return user.Users{}
	}

	for i := range users {
		log.Println("getting notes... ", users[i].Id)
		notes, err := note.GetNotes(string(users[i].Id))

		if err != nil {
			log.Println("Error getting latest ", err.Error())
			break
		}
		if len(notes) > 0 {
			log.Println("order notes... ", users[i].Id)
			users[i].Notes = notes.OldestFirst()
		}
	}
	return users
}

func NewFollowUp(s *Session) *FollowUp {
	f := &FollowUp{
		Session: s,
		Cron:    cron.New(),
	}
	f.setup([]Task{&FollowupTask{}, &CheckInTask{}, &LearnFromTask{}})
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

func (this *FollowupTask) run(fu *FollowUp) func() {
	return func() {
		log.Println("Running `Help me` ....")

		users := LoadUsersAndNotes()

		for i := range users {
			log.Println("needs help...", users[i].Id)
			help := users[i].NeedsHelp()

			if len(help) > 0 {
				fu.send(
					users[i].ToBotUser(),
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
		log.Println("checking in ....")

		users := LoadUsersAndNotes()

		for i := range users {

			log.Println("check in... ", users[i].Id)

			keyboard := Keyboard{}
			keyboard = append(keyboard, []string{"/say all good thanks"}, []string{"/chat sounds like good idea"})

			fu.sendWithKeyboard(
				users[i].ToBotUser(),
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
		log.Println("learning...")

		users := LoadUsersAndNotes()

		for i := range users {

			log.Println("learning about... ", users[i].Id)
			pos := users[i].LearnAbout(check_for_days)
			log.Println("learnt they are positive=", pos)
			keyboard := Keyboard{}

			if !pos {

				keyboard = append(keyboard, []string{"/say yeah things aren't going well"}, []string{"/say actually it is going well"})

				fu.sendWithKeyboard(
					users[i].ToBotUser(),
					"Hey, looks like things might not be going as well as you would like?",
					keyboard,
				)
				break
			}

			keyboard = append(keyboard, []string{"/say yeah I am doing well!"}, []string{"/say actually its not going well"})

			fu.sendWithKeyboard(
				users[i].ToBotUser(),
				"Hey, it looks like your doing well!!",
				keyboard,
			)

		}
	}
}

func (this *LearnFromTask) spec() string {
	//7am on MON,WED,FRI,SUN
	return "0 0 6 * * MON,WED,FRI,SUN"
}
