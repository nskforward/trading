package finam

import (
	"context"
	"crypto/tls"
	"fmt"
	"os"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/auth"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	cfg               *ClientConfig
	conn              *grpc.ClientConn
	token             string
	authService       auth.AuthServiceClient
	orderService      orders.OrdersServiceClient
	assetService      assets.AssetsServiceClient
	accountService    accounts.AccountsServiceClient
	marketDataService marketdata.MarketDataServiceClient
}

func NewClient(cfg *ClientConfig) *Client {
	if cfg == nil {
		cfg = &ClientConfig{}
	}
	return &Client{
		cfg: cfg.setDefaults(),
	}
}

func (c *Client) GetMarketDataService() (marketdata.MarketDataServiceClient, error) {
	if c.marketDataService != nil {
		return c.marketDataService, nil
	}
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	c.marketDataService = marketdata.NewMarketDataServiceClient(conn)
	return c.marketDataService, nil
}

func (c *Client) GetAccountService() (accounts.AccountsServiceClient, error) {
	if c.accountService != nil {
		return c.accountService, nil
	}
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	c.accountService = accounts.NewAccountsServiceClient(conn)
	return c.accountService, nil
}

func (c *Client) GetAssetService() (assets.AssetsServiceClient, error) {
	if c.assetService != nil {
		return c.assetService, nil
	}
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	c.assetService = assets.NewAssetsServiceClient(conn)
	return c.assetService, nil
}

func (c *Client) GetOrderService() (orders.OrdersServiceClient, error) {
	if c.orderService != nil {
		return c.orderService, nil
	}
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	c.orderService = orders.NewOrdersServiceClient(conn)
	return c.orderService, nil
}

func (c *Client) GetAuthService() (auth.AuthServiceClient, error) {
	if c.authService != nil {
		return c.authService, nil
	}
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	c.authService = auth.NewAuthServiceClient(conn)
	return c.authService, nil
}

func (c *Client) authContext(baseContext context.Context) (context.Context, error) {
	token, err := c.getToken()
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(baseContext, "Authorization", token), nil
}

func (c *Client) getToken() (string, error) {
	if c.token != "" {
		return c.token, nil
	}

	secret := os.Getenv("FINAM_SECRET")
	if secret == "" {
		return "", fmt.Errorf("env 'FINAM_SECRET' must be set")
	}

	authService, err := c.GetAuthService()
	if err != nil {
		return "", err
	}

	resp, err := authService.Auth(context.Background(), &auth.AuthRequest{Secret: secret})
	if err != nil {
		return "", fmt.Errorf("authorization failed: %w", err)
	}

	c.token = resp.GetToken()

	return c.token, nil
}

func (c *Client) getConn() (*grpc.ClientConn, error) {
	if c.conn != nil {
		return c.conn, nil
	}

	conn, err := grpc.NewClient(c.cfg.Addr, grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12})))
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}

	c.conn = conn

	return c.conn, nil
}

func (c *Client) allowTrading() error {
	if os.Getenv("ALLOW_TRADING") == "1" {
		return nil
	}
	return fmt.Errorf("trading not allowed")
}
