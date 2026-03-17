package types

import (
	"iter"
)

type Broker interface {
	GetOrders() ([]Order, error)
	SubscribeOrders() (iter.Seq[Order], error)

	GetPositions() ([]Position, error)
	SubscribeMyTrades() (iter.Seq[Position], error)

	PlaceLimitOrder(symbol string, price, size float64, pricePrec, sizePrec int) (Order, error)
	PlaceMarketOrder(symbol string, size float64, sizePrec int) (Order, error)
	CancelOrder(id string) error

	SubscribeMarketData(symbols []string) (iter.Seq[Quote], error)
	GetSchedule(symbol string) (*Schedule, error)
}
