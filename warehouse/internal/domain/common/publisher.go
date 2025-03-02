package common

type Publisher interface {
	Publish(event Event) error
}
