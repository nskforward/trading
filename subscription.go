package trading

import (
	"fmt"
	"log/slog"

	"github.com/nskforward/trading/types"
)

type Subscription struct {
	strategies      []types.Strategy
	scheduleStore   *ScheduleStore
	positionStore   *PositionStore
	limitOrderStore *LimitOrderStore
}

func NewSubscription(scheduleStore *ScheduleStore, positionStore *PositionStore, limitOrderStore *LimitOrderStore) *Subscription {
	return &Subscription{
		scheduleStore:   scheduleStore,
		limitOrderStore: limitOrderStore,
		positionStore:   positionStore,
		strategies:      make([]types.Strategy, 0, 16),
	}
}

func (s *Subscription) AddStrategy(strategy types.Strategy) error {
	for _, existing := range s.strategies {
		if existing.ID() == strategy.ID() {
			return fmt.Errorf("strategy '%s' already added before", strategy.ID())
		}
	}
	s.strategies = append(s.strategies, strategy)
	return nil
}

func (s *Subscription) Broadcast(broker types.Broker, quote types.Quote) error {
	session, err := s.scheduleStore.CurrentSession(quote.Symbol)
	if err != nil {
		return err
	}

	var position *types.Position
	loadedPosition := s.positionStore.Get(quote.Symbol)
	if loadedPosition != nil {
		position = loadedPosition
	}

	for _, strategy := range s.strategies {
		err := strategy.OnEvent(broker, quote, session, position, s.limitOrderStore.Get(quote.Symbol))
		if err != nil {
			return fmt.Errorf("an error occurred during strategy (%s) handling of quote: %w", strategy.ID(), err)
		}
	}
	return nil
}

func (s *Subscription) Init(broker types.Broker) error {
	for _, strategy := range s.strategies {
		err := strategy.Init(broker)
		if err != nil {
			return fmt.Errorf("strategy '%s' init failed: %w", strategy.ID(), err)
		}
		slog.Debug("successfully initialized strategy", "id", strategy.ID())
	}
	return nil
}
