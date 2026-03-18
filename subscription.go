package trading

import (
	"fmt"
	"log/slog"

	"github.com/nskforward/trading/types"
)

type Subscription struct {
	strategies      []*SubscribedStrategy
	scheduleStore   *ScheduleStore
	positionStore   *PositionStore
	limitOrderStore *LimitOrderStore
	assetStore      *AssetStore
}

func NewSubscription(scheduleStore *ScheduleStore, positionStore *PositionStore, limitOrderStore *LimitOrderStore, assetStore *AssetStore) *Subscription {
	return &Subscription{
		scheduleStore:   scheduleStore,
		limitOrderStore: limitOrderStore,
		positionStore:   positionStore,
		strategies:      make([]*SubscribedStrategy, 0, 16),
		assetStore:      assetStore,
	}
}

func (s *Subscription) AddStrategy(strategy types.Strategy) error {
	for _, v := range s.strategies {
		if v.strategy.ID() == strategy.ID() {
			return fmt.Errorf("strategy '%s' already added before", strategy.ID())
		}
	}
	s.strategies = append(s.strategies, NewSubscribedStrategy(s, strategy))
	return nil
}

func (s *Subscription) Broadcast(quote types.Quote) error {
	for _, v := range s.strategies {
		v.Enqueue(quote)
	}
	return nil
}

func (s *Subscription) Event(quote types.Quote) (types.Event, error) {
	session, err := s.scheduleStore.CurrentSession(quote.Symbol)
	if err != nil {
		return types.Event{}, err
	}

	asset, err := s.assetStore.Get(quote.Symbol)
	if err != nil {
		return types.Event{}, err
	}

	return types.Event{
		Quote:    quote,
		Asset:    asset,
		Session:  session,
		Position: s.positionStore.Get(quote.Symbol),
		Orders:   s.limitOrderStore.Get(quote.Symbol),
	}, nil
}

func (s *Subscription) Init() error {
	for _, v := range s.strategies {
		err := v.strategy.Init()
		if err != nil {
			return fmt.Errorf("strategy '%s' init failed: %w", v.strategy.ID(), err)
		}
		slog.Debug("successfully initialized strategy", "id", v.strategy.ID())
	}
	return nil
}
