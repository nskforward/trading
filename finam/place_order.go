package finam

import (
	"context"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
)

func (c *Client) PlaceOrder(order *orders.Order) (*orders.OrderState, error) {
	err := c.allowTrading()
	if err != nil {
		return nil, err
	}

	orderService, err := c.GetOrderService()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, err = c.authContext(ctx)

	state, err := orderService.PlaceOrder(ctx, order)
	if err != nil {
		limitPrice := "0"
		if order.LimitPrice != nil {
			limitPrice = order.LimitPrice.Value
		}
		return nil, fmt.Errorf("place order failed (type: %s, account_id: %s, symbol: %s, quantity: %s, side: %s, limit_price: %s): %w", order.Type.String(), order.AccountId, order.Symbol, order.Quantity.Value, order.Side.String(), limitPrice, err)
	}

	return state, nil
}
