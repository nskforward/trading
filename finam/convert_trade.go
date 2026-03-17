package finam

import (
	v1 "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1"
	"github.com/nskforward/trading/types"
)

func convertAccountTrade(in *v1.AccountTrade) types.Trade {
	size := convertDecimal(in.Size)
	if in.Side == v1.Side_SIDE_SELL {
		size = -size
	}

	return types.Trade{
		Time:   in.Timestamp.AsTime(),
		Symbol: in.Symbol,
		Price:  convertDecimal(in.Price),
		Size:   size,
	}
}
