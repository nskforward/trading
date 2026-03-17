package types

import "time"

type Trade struct {
	Time   time.Time
	Symbol string
	Price  float64
	Size   float64
}
