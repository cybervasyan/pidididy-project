package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
)

func (r *repository) Create(_ context.Context, req *model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	saved := *req
	r.orders[req.OrderUUID] = &saved

	return nil
}
