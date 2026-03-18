package types

type Strategy interface {
	ID() string
	Init() error
	Broker() Broker
	Symbols() []string
	OnEvent(Event) error
}

type Event struct {
	Quote    Quote
	Asset    Asset
	Session  Session
	Position *Position
	Orders   []Order
}
