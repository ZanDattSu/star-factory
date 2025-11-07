package v1

import (
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

type client struct {
	genClient inventoryV1.InventoryServiceClient
}

func NewClient(genClient inventoryV1.InventoryServiceClient) *client {
	return &client{genClient: genClient}
}
