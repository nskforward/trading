package finam

import (
	"fmt"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
	"github.com/nskforward/trading/types"
)

func convertAsset(in *assets.GetAssetResponse) types.Asset {

	fmt.Println(in.String())

	return types.Asset{
		Symbol:    in.Ticker,
		PricePrec: int(in.Decimals),
		LotSize:   convertDecimal(in.LotSize),
		MinStep:   float64(in.MinStep),
		Currency:  in.QuoteCurrency,
		Type:      in.Type,
	}
}
