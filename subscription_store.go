package trading

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/nskforward/trading/types"
)

type SubscriptionStore struct {
	broker          types.Broker
	subscriptions   map[string]*Subscription
	scheduleStore   *ScheduleStore
	positionStore   *PositionStore
	limitOrderStore *LimitOrderStore
	symbols         []string
}

func NewSubscriptionStore(broker types.Broker) *SubscriptionStore {
	return &SubscriptionStore{
		broker:          broker,
		subscriptions:   make(map[string]*Subscription),
		scheduleStore:   NewScheduleStore(broker),
		positionStore:   NewPositionStore(broker),
		limitOrderStore: NewLimitOrderStore(broker),
		symbols:         make([]string, 0, 32),
	}
}

func (store *SubscriptionStore) Init() error {
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
		if !slices.Contains(store.symbols, symbol) {
			store.symbols = append(store.symbols, symbol)
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

	stream, err := store.broker.SubscribeMarketData(store.symbols)
	if err != nil {
		return err
	}

	slog.Debug("successfully subscribed for market data", "symbols", strings.Join(store.symbols, ", "))

	for q := range stream {
		sub := store.get(q.Symbol)
		if sub == nil {
			continue
		}
		err := sub.Broadcast(store.broker, q)
		if err != nil {
			slog.Error(err.Error())
		}
	}
	return nil
}

func (store *SubscriptionStore) get(symbol string) *Subscription {
	s, ok := store.subscriptions[symbol]
	if ok {
		return s
	}
	return nil
}

func (store *SubscriptionStore) getOrCreate(symbol string) *Subscription {
	s, ok := store.subscriptions[symbol]
	if ok {
		return s
	}
	s = NewSubscription(store.scheduleStore, store.positionStore, store.limitOrderStore)
	store.subscriptions[symbol] = s
	return s
}
