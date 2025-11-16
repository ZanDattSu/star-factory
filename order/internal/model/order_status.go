package model

import "fmt"

type OrderStatus string

const (
	OrderStatusUNSPECIFIED    OrderStatus = "UNSPECIFIED"
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
	OrderStatusASSEMBLED      OrderStatus = "ASSEMBLED"
)

var orderStatusToID = map[OrderStatus]int{
	OrderStatusUNSPECIFIED:    1,
	OrderStatusPENDINGPAYMENT: 2,
	OrderStatusPAID:           3,
	OrderStatusCANCELLED:      4,
	OrderStatusASSEMBLED:      5,
}

var idToOrderStatus = map[int]OrderStatus{
	1: OrderStatusUNSPECIFIED,
	2: OrderStatusPENDINGPAYMENT,
	3: OrderStatusPAID,
	4: OrderStatusCANCELLED,
	5: OrderStatusASSEMBLED,
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
		return OrderStatusUNSPECIFIED, fmt.Errorf("unknown order status id: %d", id)
	}
	return s, nil
}
