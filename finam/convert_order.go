package finam

import (
	v1 "github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1"
	"github.com/FinamWeb/finam-trade-api/go/grpc/tradeapi/v1/orders"
	"github.com/nskforward/trading/types"
)

func convertOrder(in *orders.OrderState) types.Order {
	size := convertDecimal(in.Order.Quantity)

	if in.Order.Side == v1.Side_SIDE_SELL {
		size = -size
	}

	return types.Order{
		ID:     in.OrderId,
		Status: convertOrderStatus(in.Status),
		Symbol: in.Order.Symbol,
		Type:   convertOrderType(in.Order.Type),
		Size:   size,
		Price:  convertDecimal(in.Order.LimitPrice),
	}
}

func convertOrderType(in orders.OrderType) types.OrderType {
	switch in {
	case orders.OrderType_ORDER_TYPE_LIMIT:
		return types.OrderTypeLimit

	case orders.OrderType_ORDER_TYPE_MARKET:
		return types.OrderTypeMarket

	default:
		return types.OrderTypeUnknown
	}
}

func convertOrderStatus(in orders.OrderStatus) types.OrderStatus {
	switch in {

	case orders.OrderStatus_ORDER_STATUS_NEW:
		return types.OrderStatusNew

	case orders.OrderStatus_ORDER_STATUS_PARTIALLY_FILLED:
		return types.OrderStatusPartiallyFilled

	case orders.OrderStatus_ORDER_STATUS_FILLED:
		return types.OrderStatusFilled

	case orders.OrderStatus_ORDER_STATUS_CANCELED:
		return types.OrderStatusCanceled

	case orders.OrderStatus_ORDER_STATUS_EXPIRED:
		return types.OrderStatusExpired

	case orders.OrderStatus_ORDER_STATUS_REJECTED:
		return types.OrderStatusRejected

	default:
		return types.OrderStatusUnknown
	}
}
