package trading

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/nskforward/trading/types"
)

type SubscriptionStore struct {
	broker          types.Broker
	scheduleStore   *ScheduleStore
	positionStore   *PositionStore
	limitOrderStore *LimitOrderStore
	marketDataStore *MarketDataStore
	symbols         []string
	subscriptions   []*Subscription
}

func NewSubscriptionStore(broker types.Broker) *SubscriptionStore {
	return &SubscriptionStore{
		broker:          broker,
		subscriptions:   make([]*Subscription, 0, 16),
		scheduleStore:   NewScheduleStore(broker),
		positionStore:   NewPositionStore(broker),
		limitOrderStore: NewLimitOrderStore(broker),
		symbols:         make([]string, 0, 32),
	}
}

func (store *SubscriptionStore) InitStrategies() error {
	if len(store.subscriptions) == 0 {
		return fmt.Errorf("no subscriptions")
	}

	for _, sub := range store.subscriptions {
		err := sub.Init(store.broker)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *SubscriptionStore) AddStrategy(strategy types.Strategy) error {
	for _, symbol := range strategy.Symbols() {
		sub := store.getOrCreate(symbol)
		err := sub.AddStrategy(strategy)
		if err != nil {
			return err
		}
	}
	return nil
}

func (store *SubscriptionStore) SubscribeAndWatch() error {

	err := store.positionStore.WatchChanges()
	if err != nil {
		return err
	}

	err = store.limitOrderStore.WatchChanges()
	if err != nil {
		return err
	}

	store.marketDataStore = NewMarketDataStore(store.broker, store.symbols)
	err = store.marketDataStore.WatchChanges()
	if err != nil {
		return err
	}

	slog.Debug("successfully subscribed for market data", "symbols", strings.Join(store.symbols, ", "))

	err = store.InitStrategies()
	if err != nil {
		return fmt.Errorf("strategies init failed: %w", err)
	}

	slog.Debug("successfully initialized strategies")

	return store.watch()
}

func (store *SubscriptionStore) watch() error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if store.marketDataStore.OnQuote(store.onQuote) == 0 {
			slog.Debug("no quotes")
		}
	}

	return nil
}

func (store *SubscriptionStore) onQuote(q types.Quote) {
	sub := store.get(q.Symbol)
	if sub == nil {
		return
	}
	err := sub.Broadcast(store.broker, q)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (store *SubscriptionStore) get(symbol string) *Subscription {
	for i, v := range store.symbols {
		if v == symbol {
			return store.subscriptions[i]
		}
	}
	return nil
}

func (store *SubscriptionStore) getOrCreate(symbol string) *Subscription {
	for i, v := range store.symbols {
		if v == symbol {
			return store.subscriptions[i]
		}
	}
	store.symbols = append(store.symbols, symbol)
	sub := NewSubscription(store.scheduleStore, store.positionStore, store.limitOrderStore)
	store.subscriptions = append(store.subscriptions, sub)
	return sub
}
