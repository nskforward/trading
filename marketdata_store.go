package trading

import (
	"log/slog"
	"sync"
	"time"

	"github.com/nskforward/trading/types"
)

type MarketDataStore struct {
	broker  types.Broker
	symbols []string
	quotes  []types.Quote
	times   []time.Time
	mx      sync.Mutex
}

func NewMarketDataStore(broker types.Broker, symbols []string) *MarketDataStore {
	return &MarketDataStore{
		broker:  broker,
		symbols: symbols,
		quotes:  initQuoteSlice(symbols),
		times:   initTimeSlice(symbols),
	}
}

func (store *MarketDataStore) OnQuote(handler func(types.Quote)) int {
	count := 0
	for i := range store.symbols {
		q, ok := store.getFreshQuote(i)
		if ok {
			handler(q)
			count++
		}
	}
	return count
}

func (store *MarketDataStore) WatchChanges() error {
	stream, err := store.broker.SubscribeMarketData(store.symbols)
	if err != nil {
		return err
	}
	go func() {
		for q := range stream {
			slog.Debug("quote from broker", "symbol", q.Symbol, "time", q.Time.Unix())
			store.setQuote(q)
		}
	}()
	return nil
}

func (store *MarketDataStore) setQuote(quote types.Quote) {
	store.mx.Lock()
	defer store.mx.Unlock()

	for i, v := range store.symbols {
		if v == quote.Symbol {
			store.quotes[i].Time = quote.Time
			store.quotes[i].Ask = quote.Ask
			store.quotes[i].Bid = quote.Bid
			slog.Debug("saved quote", "quote_time", store.quotes[i].Time.Unix(), "cache_time", store.times[i].Unix())
			return
		}
	}

	slog.Warn("found quote with unsubscribed symbol", "symbol", quote.Symbol)
}

func (store *MarketDataStore) getFreshQuote(index int) (types.Quote, bool) {
	store.mx.Lock()
	defer store.mx.Unlock()

	q := store.quotes[index]
	if q.Time.After(store.times[index]) {
		store.times[index] = q.Time
		return q, true
	}

	return q, false
}

func initQuoteSlice(symbols []string) []types.Quote {
	quotes := make([]types.Quote, len(symbols))
	for i, symbol := range symbols {
		quotes[i].Symbol = symbol
	}
	return quotes
}

func initTimeSlice(symbols []string) []time.Time {
	times := make([]time.Time, len(symbols))
	for i := range symbols {
		times[i] = time.Now()
	}
	return times
}
