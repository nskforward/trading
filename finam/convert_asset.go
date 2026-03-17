package finam

import (
	"fmt"
	"math"

	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/assets"
	"github.com/nskforward/trading/types"
)

func convertAsset(in *assets.GetAssetResponse) types.Asset {
	return types.Asset{
		Symbol:    fmt.Sprintf("%s@%s", in.Ticker, in.Mic),
		PricePrec: int(in.Decimals),
		LotSize:   convertDecimal(in.LotSize),
		MinStep:   float64(in.MinStep) / math.Pow(10, float64(in.Decimals)),
		Currency:  convertCurrency(in.QuoteCurrency),
		Type:      convertAssetType(in.Type),
	}
}

func convertCurrency(in string) types.Currency {
	switch in {
	case "RUB":
		return types.RUB
	default:
		panic(fmt.Errorf("unknown currency: %s", in))
	}
}

func convertAssetType(in string) types.AssetType {
	switch in {
	case "EQUITIES":
		return types.AssetEquities
	default:
		panic(fmt.Errorf("unknown asset type: %s", in))
	}
}
