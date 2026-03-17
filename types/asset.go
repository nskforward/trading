package types

type Asset struct {
	Symbol    string
	PricePrec int
	LotSize   float64
	MinStep   float64
	Currency  string
	Type      string
}
