package finam

import (
	"context"
	"fmt"
	"iter"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
)

func (c *Client) SubscribeAccountInfo(ctx context.Context, accountID string) (iter.Seq[*accounts.GetAccountResponse], error) {
	accountService, err := c.GetAccountService()
	if err != nil {
		return nil, err
	}

	req := &accounts.GetAccountRequest{
		AccountId: accountID,
	}

	ctx, err = c.authContext(ctx)

	stream, err := accountService.SubscribeAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("subscribe account failed: %w", err)
	}

	iterator := func(yield func(*accounts.GetAccountResponse) bool) {
		for {
			resp, err := stream.Recv()
			if err != nil {
				return
			}

			if !yield(resp) {
				return
			}
		}
	}

	return iterator, nil
}
