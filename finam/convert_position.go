package finam

import (
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/accounts"
	"github.com/nskforward/trading/types"
)

func convertPosition(in *accounts.Position) types.Position {
	return types.Position{
		Symbol: in.Symbol,
		Size:   convertDecimal(in.Quantity),
		Price:  convertDecimal(in.AveragePrice),
	}
}
