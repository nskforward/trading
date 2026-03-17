package trading

import (
	"log/slog"
	"maps"
	"sync"
	"time"

	"github.com/nskforward/trading/types"
)

type PositionStore struct {
	broker     types.Broker
	positions  map[string]*types.Position
	mx         sync.RWMutex
	started    time.Time
	assetStore *AssetStore
}

func NewPositionStore(broker types.Broker, assetStore *AssetStore) *PositionStore {
	return &PositionStore{
		broker:     broker,
		positions:  make(map[string]*types.Position),
		started:    time.Now(),
		assetStore: assetStore,
	}
}

func (store *PositionStore) Get(symbol string) *types.Position {
	store.mx.RLock()
	defer store.mx.RUnlock()

	pos, ok := store.positions[symbol]
	if ok {
		return pos
	}
	return nil
}

func (store *PositionStore) WatchChanges() error {
	err := store.refresh()
	if err != nil {
		return err
	}

	stream, err := store.broker.SubscribeMyTrades()
	if err != nil {
		return nil
	}
	go func() {
		for trade := range stream {
			if trade.Time.Before(store.started) {
				continue
			}

			slog.Debug("found trade", "symbol", trade.Symbol, "size", trade.Size, "price", trade.Price)

			pos, err := store.update(trade)
			if err != nil {
				slog.Error("position chnage failed", "error", err.Error())
				continue
			}

			if pos.Size == 0 {
				store.delete(pos.Symbol)
			}
		}
	}()
	return nil
}

func (store *PositionStore) delete(symbol string) {
	store.mx.Lock()
	defer store.mx.Unlock()
	delete(store.positions, symbol)
}

func (store *PositionStore) update(trade types.Trade) (*types.Position, error) {
	pos := store.Get(trade.Symbol)
	if pos == nil {
		pos = &types.Position{}
	}

	asset, err := store.assetStore.Get(trade.Symbol)
	if err != nil {
		return nil, err
	}

	pos.Merge(asset, trade)

	if pos.Symbol == "" {
		pos.Symbol = trade.Symbol
		store.set(pos)
	}

	return pos, nil
}

func (store *PositionStore) set(position *types.Position) {
	store.mx.Lock()
	defer store.mx.Unlock()
	store.positions[position.Symbol] = position
}

func (store *PositionStore) refresh() error {
	newValues, err := store.broker.GetPositions()
	if err != nil {
		return err
	}

	store.mx.Lock()
	defer store.mx.Unlock()

	maps.DeleteFunc(store.positions, func(k string, v *types.Position) bool { return true })

	for _, newValue := range newValues {
		store.positions[newValue.Symbol] = &newValue
	}

	return nil
}
