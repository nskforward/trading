package finam

import (
	"context"
	"fmt"
	"iter"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
)

func (c *Client) SubscribeQuotes(ctx context.Context, symbols []string) (iter.Seq[*marketdata.Quote], error) {
	marketDataService, err := c.GetMarketDataService()
	if err != nil {
		return nil, err
	}

	req := &marketdata.SubscribeQuoteRequest{
		Symbols: symbols,
	}

	ctx, err = c.authContext(ctx)

	stream, err := marketDataService.SubscribeQuote(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("subscribe quotes failed: %w", err)
	}

	iterator := func(yield func(*marketdata.Quote) bool) {
		for {
			resp, err := stream.Recv()
			if err != nil {
				return
			}
			for _, order := range resp.Quote {
				if !yield(order) {
					return
				}
			}
		}
	}

	return iterator, nil
}
