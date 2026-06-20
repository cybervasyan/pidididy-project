package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, req *model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[req.OrderUUID]

	if !ok {
		return model.Order{}, model.ErrorOrderDoesntExist
	}

	if order.Status != model.OrderStatusPENDINGPAYMENT {
		return model.Order{}, model.ErrorOrderNotInPending
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = req.PaymentMethod
	order.TransactionUUID = req.TransactionUUID
	return *order, nil
}
