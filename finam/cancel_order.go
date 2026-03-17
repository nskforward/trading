package finam

import (
	"context"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
)

func (c *Client) CancelOrder(accountID, orderID string) (*orders.OrderState, error) {
	err := c.allowTrading()
	if err != nil {
		return nil, err
	}

	orderService, err := c.GetOrderService()
	if err != nil {
		return nil, err
	}

	req := &orders.CancelOrderRequest{
		AccountId: accountID,
		OrderId:   orderID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, err = c.authContext(ctx)

	state, err := orderService.CancelOrder(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("cancel order failed: %w", err)
	}

	return state, nil
}
