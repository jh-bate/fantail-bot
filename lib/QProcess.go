package lib

import "github.com/jh-bate/fantail-bot/Godeps/_workspace/src/github.com/tucnak/telebot"

const (
	default_script = "default"
	stickers_chat  = "chat"
)

type (
	Info struct {
		App       []string `json:"appInfo"`
		Reminders []string `json:"remindersInfo"`
		Chat      []string `json:"chatInfo"`
		Said      []string `json:"saidInfo"`
	}

	QProcess struct {
		s *session
		/*lang struct {
			Questions `json:"QandA"`
		}
		info     *Info
		next     *Question
		in       *Incoming
		lastTime Notes
		sLib     Stickers*/
	}
)

func NewQProcess(b *telebot.Bot, s *Storage) *QProcess {
	q := &QProcess{s: newSession(b, s)}
	//q.loadInfo()
	return q
}

func (this *QProcess) Run(input <-chan telebot.Message) {
	for msg := range input {

		//this.in = newIncoming(msg).getAction(this.s)
		//this.s.User = this.in.sender()

		action := newIncoming(msg).getAction(this.s)
		action.doFirst()
		action.loadQuestions()
		action.chat(action.findNext())

		/*if this.in.isSticker() {
			log.Println("incoming sticker", msg.Sticker.FileID)
			if s := this.sLib.FindSticker(msg.Sticker.FileID); s != nil {
				this.loadScript(stickers_chat)
				this.
					nextStickerQ(s).
					andChat()
			}
		} else {
			this.
				quickWinFirst().
				determineScript().
				nextQ().
				andChat()
		}*/

	}
}

/*
func (this *QProcess) quickWinFirst() *QProcess {
	this.in.getAction(this.s).doFirst()
	return this
}

func (this *QProcess) loadInfo() {
	file, err := os.Open("./config/fantail.json")
	if err != nil {
		log.Panic("could not load App info", err.Error())
	}
	err = json.NewDecoder(file).Decode(&this.info)
	if err != nil {
		log.Panic("could not decode App info ", err.Error())
	}
}

func (this *QProcess) loadScript(scriptName string) {
	file, err := os.Open(fmt.Sprintf("./config/%s.json", scriptName))
	if err != nil {
		log.Panic("could not load QandA language file ", err.Error())
	}
	err = json.NewDecoder(file).Decode(&this.lang)
	if err != nil {
		log.Panic("could not decode QandA ", err.Error())
	}
}

func (this *QProcess) determineScript() *QProcess {

	if this.in.submissonMatches(remind_action, say_action) {
		log.Println("load default script after submisson")
		this.loadScript(default_script)
	} else if this.in.cmdMatches(chat_action, say_action, remind_action) {
		log.Println("load command script", this.in.getCmd())
		this.loadScript(this.in.getCmd())
	}
	return this
}

func (this *QProcess) nextQ() *QProcess {
	this.next = nil

	if this.in.cmdMatches(chat_action, say_action, remind_action) {
		this.next = this.lang.Questions.First()
		return this
	} else {
		if nxt, sv := this.lang.Questions.next(this.in.msg.Text); sv {
			this.next = nxt
			this.s.save(this.in.getNote(chat_action, this.next.RelatesTo.SaveTag))
			return this
		} else {
			this.next = nxt
			return this
		}
	}
}

func (this *QProcess) nextStickerQ(s *Sticker) *QProcess {
	this.next = nil

	if nxt, sv := this.lang.Questions.nextFrom(s.Ids...); sv {
		this.s.save(this.in.getNote(chat_tag))
		this.next = nxt
		return this
	} else {
		this.next = nxt
		return this
	}
}

func (this *QProcess) andChat() {
	if this.next != nil {
		this.s.send(this.next.Context...)
		this.s.sendWithKeyboard(this.next.QuestionText, this.next.makeKeyboard())
	}
	return
}
*/
