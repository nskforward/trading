package finam

import (
	"context"
	"fmt"
	"iter"
	"os"
	"strconv"

	v1 "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
	"github.com/nskforward/trading/types"
	"google.golang.org/genproto/googleapis/type/decimal"
)

type Broker struct {
	client     *Client
	accountID  string
	quoteCache *quoteCache
}

func NewBroker() *Broker {
	accountID := os.Getenv("FINAM_ACCOUNT")
	if accountID == "" {
		panic("env 'FINAM_ACCOUNT' is not set")
	}

	return &Broker{
		client:     NewClient(nil),
		accountID:  accountID,
		quoteCache: newQuoteCache(),
	}
}

func (b *Broker) GetPositions() ([]types.Position, error) {
	resp, err := b.client.GetAccount(b.accountID)
	if err != nil {
		return nil, fmt.Errorf("cannot get account info: %w", err)
	}
	result := make([]types.Position, 0, len(resp.Positions))
	for _, pos := range resp.Positions {
		result = append(result, convertPosition(pos))
	}

	return result, nil
}

func (b *Broker) GetOrders() ([]types.Order, error) {
	resp, err := b.client.GetOrders(b.accountID)
	if err != nil {
		return nil, err
	}

	list := make([]types.Order, 0, len(resp))

	for _, state := range resp {
		list = append(list, convertOrder(state))
	}

	return list, nil
}

func (b *Broker) GetSchedule(symbol string) (*types.Schedule, error) {
	slots, err := b.client.GetSchedule(symbol)
	if err != nil {
		return nil, err
	}
	schedule := types.Schedule{
		Slots: make([]types.Session, 0, len(slots)),
	}
	for _, slot := range slots {
		schedule.Slots = append(schedule.Slots, types.Session{
			Type:  convertScheduleSlotType(slot.Type),
			Start: slot.Interval.StartTime.Seconds,
			End:   slot.Interval.EndTime.Seconds,
		})
	}
	return &schedule, nil
}

func (b *Broker) CancelOrder(id string) error {
	_, err := b.client.CancelOrder(b.accountID, id)
	return err
}

func (b *Broker) PlaceLimitOrder(symbol string, price, size float64, pricePrec, sizePrec int) (types.Order, error) {
	side := v1.Side_SIDE_BUY
	if size < 0 {
		side = v1.Side_SIDE_SELL
		size = -size
	}

	order := &orders.Order{
		AccountId: b.accountID,
		Symbol:    symbol,
		Quantity: &decimal.Decimal{
			Value: strconv.FormatFloat(size, 'f', sizePrec, 64),
		},
		Side:        side,
		Type:        orders.OrderType_ORDER_TYPE_LIMIT,
		TimeInForce: orders.TimeInForce_TIME_IN_FORCE_DAY,
		LimitPrice: &decimal.Decimal{
			Value: strconv.FormatFloat(price, 'f', pricePrec, 64),
		},
	}

	state, err := b.client.PlaceOrder(order)
	if err != nil {
		return types.Order{}, err
	}

	return convertOrder(state), nil
}

func (b *Broker) PlaceMarketOrder(symbol string, size float64, sizePrec int) (types.Order, error) {
	side := v1.Side_SIDE_BUY
	if size < 0 {
		side = v1.Side_SIDE_SELL
		size = -size
	}

	order := &orders.Order{
		Type:      orders.OrderType_ORDER_TYPE_MARKET,
		AccountId: b.accountID,
		Symbol:    symbol,
		Quantity: &decimal.Decimal{
			Value: strconv.FormatFloat(size, 'f', sizePrec, 64),
		},
		Side: side,
	}

	state, err := b.client.PlaceOrder(order)
	if err != nil {
		return types.Order{}, err
	}

	return convertOrder(state), nil
}

func (b *Broker) SubscribeMyTrades() (iter.Seq[types.Position], error) {
	stream, err := b.client.SubscribeMyTrades(context.Background(), b.accountID)
	if err != nil {
		return nil, err
	}

	iterator := func(yield func(types.Position) bool) {
		for in := range stream {
			if !yield(convertPositionFromTrade(in)) {
				return
			}
		}
	}
	return iterator, nil
}

func (b *Broker) SubscribeOrders() (iter.Seq[types.Order], error) {
	stream, err := b.client.SubscribeOrders(context.Background(), b.accountID)
	if err != nil {
		return nil, err
	}
	iterator := func(yield func(types.Order) bool) {
		for in := range stream {
			if !yield(convertOrder(in)) {
				return
			}
		}
	}
	return iterator, nil
}

func (b *Broker) SubscribeMarketData(symbols []string) (iter.Seq[types.Quote], error) {
	stream, err := b.client.SubscribeQuotes(context.Background(), symbols)
	if err != nil {
		return nil, err
	}
	return func(yield func(types.Quote) bool) {
		for quote := range stream {
			in := types.Quote{
				Symbol: quote.Symbol,
				Ask:    convertDecimal(quote.Ask),
				Bid:    convertDecimal(quote.Bid),
			}

			out := b.quoteCache.Get(in)
			if out != nil {
				if !yield(*out) {
					return
				}
			}
		}
	}, nil
}

func convertDecimal(in *decimal.Decimal) float64 {
	if in == nil {
		return 0
	}
	num, err := strconv.ParseFloat(in.Value, 64)
	if err != nil {
		return 0
	}
	return num
}

func convertScheduleSlotType(in string) types.SessionType {
	switch in {
	case "EARLY_TRADING", "OPENING_AUCTION":
		return types.SessionPremarket
	case "CORE_TRADING":
		return types.SessionMain
	case "LATE_TRADING", "CLOSING_AUCTION":
		return types.SessionPostmarket
	default:
		return types.SessionClosed
	}
}
