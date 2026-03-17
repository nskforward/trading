package finam

import (
	"context"
	"fmt"
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
)

func (c *Client) GetSchedule(symbol string) ([]*assets.ScheduleResponse_Sessions, error) {
	assetService, err := c.GetAssetService()
	if err != nil {
		return nil, err
	}

	req := &assets.ScheduleRequest{
		Symbol: symbol,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ctx, err = c.authContext(ctx)

	resp, err := assetService.Schedule(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("get schedule failed: %w", err)
	}

	return resp.Sessions, nil
}
