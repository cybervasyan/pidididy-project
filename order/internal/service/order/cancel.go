package order

import (
	"context"

	baseModel "github.com/cybervasyan/pdididy-project/order/internal/model"
	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
	"github.com/google/uuid"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if order.Status != model.OrderStatusPENDINGPAYMENT {
		return baseModel.ErrorOrderNotInPending
	}

	err = s.orderRepo.Cancel(ctx, orderUUID)
	if err != nil {
		return err
	}

	return nil
}
