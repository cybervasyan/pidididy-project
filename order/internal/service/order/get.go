package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/model"
	"github.com/cybervasyan/pdididy-project/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (s *service) GetOrderByUuid(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		return model.Order{}, err
	}

	return converter.ToServiceModel(order), nil
}
