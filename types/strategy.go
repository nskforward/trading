package types

type Strategy interface {
	ID() string
	Symbols() []string
	Init(Broker) error

	// OnEvent returns only critical error that trigger to process exit
	OnEvent(EventContext) error
}

type EventContext struct {
	Quote    Quote
	Broker   Broker
	Session  Session
	Position *Position
	Orders   []Order
}
