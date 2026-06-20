package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) CreateOrder(ctx context.Context, req *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[req.OrderUUID] = &model.Order{
		OrderUUID: req.OrderUUID,
		UserUUID:  req.UserUUID,
		PartUuids: req.PartUuids,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	return nil
}
