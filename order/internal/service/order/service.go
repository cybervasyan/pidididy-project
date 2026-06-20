package order

import (
	"github.com/cybervasyan/pdididy-project/order/internal/client/grpc"
	"github.com/cybervasyan/pdididy-project/order/internal/repository"
)

type service struct {
	orderRepo       repository.Repository
	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewOrderService(orderRepo repository.Repository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepo:       orderRepo,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
