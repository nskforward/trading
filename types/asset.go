package types

type Asset struct {
	Symbol    string
	PricePrec int
	LotSize   float64
	MinStep   float64
	Currency  Currency
	Type      AssetType
}

type AssetType string

const (
	AssetEquities AssetType = "equities"
)

type Currency string

const (
	RUB Currency = "rub"
)
