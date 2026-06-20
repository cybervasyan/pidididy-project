package order

import (
	"github.com/cybervasyan/pdididy-project/order/internal/client/grpc"
	"github.com/cybervasyan/pdididy-project/order/internal/repository"
)

type service struct {
	orderRepo       repository.Repository
	inventoryClient grpc.InventoryClient
}

func NewOrderService(orderRepo repository.Repository, inventoryClient grpc.InventoryClient) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
	}
}
