package bot

type Bot interface {
	Listen(subscription chan<- *Payload)
	SendMessage(recipientId, message string) error
}
