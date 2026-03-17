package finam

import (
	"context"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
)

func (c *Client) GetOrders(accountID string) ([]*orders.OrderState, error) {
	orderService, err := c.GetOrderService()
	if err != nil {
		return nil, err
	}

	req := &orders.OrdersRequest{
		AccountId: accountID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, err = c.authContext(ctx)

	state, err := orderService.GetOrders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get orders failed: %w", err)
	}

	return state.GetOrders(), nil
}
