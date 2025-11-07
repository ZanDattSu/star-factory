package model

import (
	"fmt"
)

type OrderNotFoundError struct {
	Code      int    `json:"code"`
	OrderUUID string `json:"order_uuid"`
}

func (e *OrderNotFoundError) Error() string {
	return fmt.Sprintf("order with UUID %q not found", e.OrderUUID)
}

func NewOrderNotFoundError(uuid string) *OrderNotFoundError {
	return &OrderNotFoundError{
		Code:      404,
		OrderUUID: uuid,
	}
}

type ConflictError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *ConflictError) Error() string {
	return e.Message
}

func NewConflictError(message string) *ConflictError {
	return &ConflictError{
		Code:    409,
		Message: message,
	}
}
