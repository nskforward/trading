package trading

import (
	"log/slog"
	"maps"
	"sync"

	"github.com/nskforward/trading/types"
)

type LimitOrderStore struct {
	broker types.Broker
	orders map[string][]types.Order
	mx     sync.RWMutex
}

func NewLimitOrderStore(broker types.Broker) *LimitOrderStore {
	return &LimitOrderStore{
		broker: broker,
		orders: make(map[string][]types.Order),
	}
}

func (store *LimitOrderStore) Get(symbol string) []types.Order {
	store.mx.RLock()
	defer store.mx.RUnlock()

	res, ok := store.orders[symbol]
	if ok {
		return res
	}
	return nil
}

func (store *LimitOrderStore) WatchChanges() error {
	err := store.refresh()
	if err != nil {
		return err
	}

	stream, err := store.broker.SubscribeOrders()
	if err != nil {
		return nil
	}
	go func() {
		for order := range stream {
			if order.Type != types.OrderTypeLimit {
				continue
			}

			slog.Debug("found order", "symbol", order.Symbol, "size", order.Size, "price", order.Price, "status", order.Status)

			store.set(order)
		}
	}()
	return nil
}

func (store *LimitOrderStore) delete(symbol string) {
	store.mx.Lock()
	defer store.mx.Unlock()
	delete(store.orders, symbol)
}

func (store *LimitOrderStore) refresh() error {
	list, err := store.broker.GetOrders()
	if err != nil {
		return err
	}

	store.clear()

	for _, order := range list {
		if order.Type == types.OrderTypeLimit {
			store.set(order)
		}
	}

	return nil
}

func (store *LimitOrderStore) clear() {
	store.mx.Lock()
	defer store.mx.Unlock()
	maps.DeleteFunc(store.orders, func(k string, v []types.Order) bool { return true })
}

func (store *LimitOrderStore) set(order types.Order) {
	store.mx.Lock()
	defer store.mx.Unlock()

	list, ok := store.orders[order.Symbol]
	if !ok {
		list = make([]types.Order, 0, 16)
	}

	updated := false

	for i, v := range list {
		if v.ID == order.ID {
			list[i] = order
			updated = true
			break
		}
	}

	if !updated {
		list = append(list, order)
	}

	filtered := list[:0]
	for _, v := range list {
		if v.Status == types.OrderStatusNew || v.Status == types.OrderStatusPartiallyFilled {
			filtered = append(filtered, v)
		}
	}

	if len(filtered) == 0 {
		delete(store.orders, order.Symbol)
	} else {
		store.orders[order.Symbol] = filtered
	}
}
