package types

type Strategy interface {
	ID() string
	Symbols() []string
	Init(Broker) error
	OnEvent(Event) error
}

type Event struct {
	Quote    Quote
	Asset    Asset
	Broker   Broker
	Session  Session
	Position *Position
	Orders   []Order
}
