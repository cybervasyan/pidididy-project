package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/model"
	"github.com/cybervasyan/pdididy-project/order/internal/repository/converter"
)

func (s *service) PayOrder(ctx context.Context, req *model.Order) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, req.OrderUUID)
	if err != nil {
		return model.Order{}, err
	}

	orderModel := converter.ToServiceModel(order)
	if orderModel.Status != model.OrderStatusPENDINGPAYMENT {
		return model.Order{}, model.ErrorOrderNotInPending
	}

	orderModel.Status = model.OrderStatusPAID
	orderModel.PaymentMethod = req.PaymentMethod
	orderModel.TransactionUUID = req.TransactionUUID

	repoOrder := converter.ToRepoModel(orderModel)
	updatedOrder, err := s.orderRepo.Update(ctx, &repoOrder)
	if err != nil {
		return model.Order{}, err
	}

	return converter.ToServiceModel(updatedOrder), nil
}
