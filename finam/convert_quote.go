package finam

import (
	"time"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/marketdata"
	"github.com/nskforward/trading/types"
)

func convertQuote(in *marketdata.Quote) types.Quote {
	return types.Quote{
		Time:   time.Now(),
		Symbol: in.Symbol,
		Ask:    convertDecimal(in.Ask),
		Bid:    convertDecimal(in.Bid),
	}
}
