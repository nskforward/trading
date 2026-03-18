package trading

import (
	"log/slog"

	"github.com/nskforward/trading/types"
)

type SubscribedStrategy struct {
	subscription *Subscription
	stream       chan types.Quote
	strategy     types.Strategy
	broker       types.Broker
}

func NewSubscribedStrategy(subscription *Subscription, strategy types.Strategy) *SubscribedStrategy {
	s := &SubscribedStrategy{
		subscription: subscription,
		stream:       make(chan types.Quote, 1),
		strategy:     strategy,
		broker:       strategy.Broker(),
	}
	go s.subscribe()
	return s
}

func (s *SubscribedStrategy) Enqueue(quote types.Quote) {
	for {
		select {
		case s.stream <- quote:
			return
		default:
			<-s.stream
		}
	}
}

func (s *SubscribedStrategy) subscribe() {
	for q := range s.stream {
		event, err := s.subscription.Event(q)
		if err != nil {
			slog.Error("cannot create event", "error", err.Error())
			continue
		}
		err = s.strategy.OnEvent(event)
		if err != nil {
			slog.Error("strategy event handling error", "strategy", s.strategy.ID(), "error", err.Error())
			continue
		}
	}
}
