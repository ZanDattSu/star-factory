package model

type Order struct {
	OrderUUID       string        `json:"order_uuid"`
	UserUUID        string        `json:"user_uuid"`
	PartUuids       []string      `json:"part_uuids"`
	TotalPrice      float64       `json:"total_price"`
	TransactionUUID *string       `json:"transaction_uuid,omitempty"`
	PaymentMethod   PaymentMethod `json:"payment_method,omitempty"`
	Status          OrderStatus   `json:"status,omitempty"`
}
