package finam

import (
	"context"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
)

func (c *Client) GetAccount(accountID string) (*accounts.GetAccountResponse, error) {
	accountService, err := c.GetAccountService()
	if err != nil {
		return nil, err
	}

	req := &accounts.GetAccountRequest{
		AccountId: accountID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, err = c.authContext(ctx)

	resp, err := accountService.GetAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get account failed: %w", err)
	}

	return resp, nil
}
