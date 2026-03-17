package finam

import (
	"context"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
)

func (c *Client) GetAsset(accountID, symbol string) (*assets.GetAssetResponse, error) {
	assetService, err := c.GetAssetService()
	if err != nil {
		return nil, err
	}

	req := &assets.GetAssetRequest{
		AccountId: accountID,
		Symbol:    symbol,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, err = c.authContext(ctx)

	resp, err := assetService.GetAsset(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get asset failed: %w", err)
	}

	return resp, nil
}
