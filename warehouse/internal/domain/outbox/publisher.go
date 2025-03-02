package outbox

type Publisher interface {
	Publish(message *Message) error
}
