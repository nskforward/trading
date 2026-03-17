package types

type Strategy interface {
	ID() string
	Symbols() []string
	Init(Broker) error

	// OnEvent returns only critical error that trigger to process exit
	OnEvent(Broker, Quote, Session, *Position, []Order) error
}
