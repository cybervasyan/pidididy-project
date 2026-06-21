package order

import (
	"context"
	"errors"

	"github.com/cybervasyan/pdididy-project/order/internal/model"
	repoConverter "github.com/cybervasyan/pdididy-project/order/internal/repository/converter"
	repoModel "github.com/cybervasyan/pdididy-project/order/internal/repository/model"
	"github.com/google/uuid"
)

func (s *service) GetOrderByUuid(ctx context.Context, orderUUID uuid.UUID) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, repoModel.ErrOrderDoesntExist) {
			return model.Order{}, model.ErrOrderDoesntExist
		}
		return model.Order{}, err
	}

	return repoConverter.ToServiceModel(order), nil
}
