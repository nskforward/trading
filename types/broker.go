package types

import (
	"iter"
)

type Broker interface {
	GetOrders() ([]Order, error)
	GetPositions() ([]Position, error)
	GetAsset(symbol string) (Asset, error)
	GetSchedule(symbol string) (*Schedule, error)

	SubscribeOrders() (iter.Seq[Order], error)
	SubscribeMyTrades() (iter.Seq[Trade], error)
	SubscribeQuotes(symbols []string) (iter.Seq[Quote], error)

	CancelOrder(id string) error
	PlaceMarketOrder(symbol string, size float64, sizePrec int) (Order, error)
	PlaceLimitOrder(symbol string, price, size float64, pricePrec, sizePrec int) (Order, error)
}
