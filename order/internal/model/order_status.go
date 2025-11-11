package model

import "fmt"

type OrderStatus string

const (
	OrderStatusUnspecified    OrderStatus = "UNSPECIFIED"
	OrderStatusPendingPayment OrderStatus = "PENDING_PAYMENT"
	OrderStatusPaid           OrderStatus = "PAID"
	OrderStatusCancelled      OrderStatus = "CANCELLED"
)

var orderStatusToID = map[OrderStatus]int{
	OrderStatusUnspecified:    1,
	OrderStatusPendingPayment: 2,
	OrderStatusPaid:           3,
	OrderStatusCancelled:      4,
}

var idToOrderStatus = map[int]OrderStatus{
	1: OrderStatusUnspecified,
	2: OrderStatusPendingPayment,
	3: OrderStatusPaid,
	4: OrderStatusCancelled,
}

func (s OrderStatus) ID() (int, error) {
	id, ok := orderStatusToID[s]
	if !ok {
		return 0, fmt.Errorf("unknown order status: %s", s)
	}
	return id, nil
}

func OrderStatusFromID(id int) (OrderStatus, error) {
	s, ok := idToOrderStatus[id]
	if !ok {
		return OrderStatusUnspecified, fmt.Errorf("unknown order status id: %d", id)
	}
	return s, nil
}
