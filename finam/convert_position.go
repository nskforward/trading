package finam

import (
	"time"

	v1 "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
	"github.com/nskforward/trading/types"
)

func convertPosition(in *accounts.Position) types.Position {
	return types.Position{
		Symbol:  in.Symbol,
		Size:    convertDecimal(in.Quantity),
		Price:   convertDecimal(in.AveragePrice),
		Updated: time.Now(),
	}
}

func convertPositionFromTrade(in *v1.AccountTrade) types.Position {
	size := convertDecimal(in.Size)

	if in.Side == v1.Side_SIDE_SELL {
		size = -size
	}

	return types.Position{
		Symbol:  in.Symbol,
		Size:    size,
		Price:   convertDecimal(in.Price),
		Updated: in.Timestamp.AsTime(),
	}
}
