package v1

import (
	"context"
	"errors"
	"log"

	"github.com/cybervasyan/pdididy-project/payment/internal/converter"
	"github.com/cybervasyan/pdididy-project/payment/internal/model"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	orderUUID, err := uuid.Parse(req.GetOrderUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "невалидный order uuid: %v", err)
	}

	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "невалидный user uuid: %v", err)
	}

	transactionUUID, err := a.paymentService.PayOrder(ctx, orderUUID, userUUID, converter.PaymentMethodToModel(req.GetPaymentMethod()))
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderUUIDRequired),
			errors.Is(err, model.ErrUserUUIDRequired),
			errors.Is(err, model.ErrPaymentMethodUnspecified):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			log.Printf("PayOrder: непредвиденная ошибка: %v", err)
			return nil, status.Error(codes.Internal, "что-то пошло не так")
		}
	}

	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID.String(),
	}, nil
}
