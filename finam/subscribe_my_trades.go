package finam

import (
	"context"
	"fmt"
	"iter"

	v1 "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
)

func (c *Client) SubscribeMyTrades(ctx context.Context, accountID string) (iter.Seq[*v1.AccountTrade], error) {
	orderService, err := c.GetOrderService()
	if err != nil {
		return nil, err
	}

	req := &orders.SubscribeTradesRequest{
		AccountId: accountID,
	}

	ctx, err = c.authContext(ctx)

	stream, err := orderService.SubscribeTrades(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("subscribe trades failed: %w", err)
	}

	iterator := func(yield func(*v1.AccountTrade) bool) {
		for {
			resp, err := stream.Recv()
			if err != nil {
				return
			}
			for _, trade := range resp.Trades {
				if !yield(trade) {
					return
				}
			}
		}
	}

	return iterator, nil
}
