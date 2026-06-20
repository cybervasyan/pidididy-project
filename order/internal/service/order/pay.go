package order

import (
	"context"

	grpcConverter "github.com/cybervasyan/pdididy-project/order/internal/client/converter"
	"github.com/cybervasyan/pdididy-project/order/internal/model"
	repoConverter "github.com/cybervasyan/pdididy-project/order/internal/repository/converter"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
)

func (s *service) PayOrder(ctx context.Context, req *model.Order) (model.Order, error) {
	order, err := s.orderRepo.Get(ctx, req.OrderUUID)
	if err != nil {
		return model.Order{}, err
	}

	orderModel := repoConverter.ToServiceModel(order)
	if orderModel.Status != model.OrderStatusPENDINGPAYMENT {
		return model.Order{}, model.ErrOrderNotInPending
	}

	paymentMethod := grpcConverter.PaymentMethodToProto(req.PaymentMethod)

	payResp, err := s.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     req.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return model.Order{}, err
	}

	parsedUUID, err := uuid.Parse(payResp.GetTransactionUuid())
	if err != nil {
		return model.Order{}, err
	}
	orderModel.Status = model.OrderStatusPAID
	orderModel.PaymentMethod = req.PaymentMethod
	orderModel.TransactionUUID = &parsedUUID

	repoOrder := repoConverter.ToRepoModel(orderModel)
	updatedOrder, err := s.orderRepo.Update(ctx, &repoOrder)
	if err != nil {
		return model.Order{}, err
	}

	return repoConverter.ToServiceModel(updatedOrder), nil
}
