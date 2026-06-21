package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Update(_ context.Context, req *model.Order) (model.Order, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.orders[req.OrderUUID]
	if !ok {
		return model.Order{}, model.ErrOrderDoesntExist
	}

	updated := *req
	r.orders[req.OrderUUID] = &updated
	return updated, nil
}
