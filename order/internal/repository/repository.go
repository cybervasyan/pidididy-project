package repository

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
	"github.com/google/uuid"
)

type Order interface {
	CreateOrder(ctx context.Context, req *model.Order) error
	CancelOrder(_ context.Context, orderUUID uuid.UUID) error
	Update(ctx context.Context, req *model.Order) (model.Order, error)
	GetOrderByUuid(_ context.Context, orderUUID uuid.UUID) (model.Order, error)
}
