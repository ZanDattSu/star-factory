package part

import (
	"github.com/ZanDattSu/star-factory/inventory/internal/service"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
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
