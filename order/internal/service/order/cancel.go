package order

import (
	"context"
	"errors"

	model "github.com/cybervasyan/pdididy-project/order/internal/model"
	repoModel "github.com/cybervasyan/pdididy-project/order/internal/repository/model"
	"github.com/google/uuid"
)

func (s *service) CancelOrder(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, repoModel.ErrOrderDoesntExist) {
			return model.ErrOrderDoesntExist
		}
		return err
	}

	if order.Status != repoModel.OrderStatusPENDINGPAYMENT {
		return model.ErrOrderNotInPending
	}

	err = s.orderRepo.Cancel(ctx, orderUUID)
	if err != nil {
		return err
	}

	return nil
}
