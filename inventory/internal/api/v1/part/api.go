package part

import (
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	"inventory/internal/service"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer
	partService service.PartService
}

func NewApi(partService service.PartService) *api {
	return &api{
		partService: partService,
	}
}
