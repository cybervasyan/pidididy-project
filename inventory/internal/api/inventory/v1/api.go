package v1

import (
	"github.com/cybervasyan/pdididy-project/inventory/internal/service"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryv1.UnimplementedInventoryServiceServer

	inventoryService service.Part
}

func NewAPI(inventoryService service.Part) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
