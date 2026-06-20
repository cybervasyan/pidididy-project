package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) CancelOrder(_ context.Context, orderUUID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.ErrorOrderDoesntExist
	}

	if order.Status != model.OrderStatusPENDINGPAYMENT {
		return model.ErrorOrderNotInPending
	}

	order.Status = model.OrderStatusCANCELLED

	return nil
}
