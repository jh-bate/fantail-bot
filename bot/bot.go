package bot

type Bot interface {
	Listen(subscription chan<- *Payload)
	SendMessage(recipientId int, message string) error
}
