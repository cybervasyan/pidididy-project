package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, req *model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.orders[req.OrderUUID]
	if !ok {
		return model.Order{}, model.ErrorOrderDoesntExist
	}

	r.orders[req.OrderUUID] = req
	return *req, nil
}
