package finam

import (
	"context"
	"fmt"
	"iter"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
)

func (c *Client) SubscribeOrders(ctx context.Context, accountID string) (iter.Seq[*orders.OrderState], error) {
	orderService, err := c.GetOrderService()
	if err != nil {
		return nil, err
	}

	req := &orders.SubscribeOrdersRequest{
		AccountId: accountID,
	}

	ctx, err = c.authContext(ctx)

	stream, err := orderService.SubscribeOrders(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("subscribe orders failed: %w", err)
	}

	iterator := func(yield func(*orders.OrderState) bool) {
		for {
			resp, err := stream.Recv()
			if err != nil {
				return
			}
			for _, order := range resp.Orders {
				if !yield(order) {
					return
				}
			}
		}
	}

	return iterator, nil
}
