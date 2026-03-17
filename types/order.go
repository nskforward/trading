package types

type Order struct {
	ID     string
	Symbol string
	Status OrderStatus
	Type   OrderType
	Size   float64
	Price  float64
}

type OrderType uint8

const (
	OrderTypeUnknown OrderType = iota
	OrderTypeMarket
	OrderTypeLimit
)

type OrderStatus string

const (
	OrderStatusUnknown         OrderStatus = "unknown"
	OrderStatusNew                         = "new"
	OrderStatusPartiallyFilled             = "partially_filled"
	OrderStatusFilled                      = "filled"
	OrderStatusCanceled                    = "canceled"
	OrderStatusRejected                    = "rejected"
	OrderStatusExpired                     = "expired"
)
