package v1

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/cybervasyan/pdididy-project/order/internal/converter"
	"github.com/cybervasyan/pdididy-project/order/internal/model"
	orderv1 "github.com/cybervasyan/pdididy-project/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (r orderv1.PayOrderRes, _ error) {
	orderModel := converter.PayRequestToModel(params.OrderUUID, req)

	order, err := a.orderService.PayOrder(ctx, &orderModel)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderDoesntExist):
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Заказ не найден",
			}, nil

		case errors.Is(err, model.ErrOrderNotInPending):
			return &orderv1.ConflictError{
				Code:    http.StatusConflict,
				Message: "Заказ не в статусе PENDING_PAYMENT",
			}, nil

		default:
			log.Printf("PayOrder: непредвиденная ошибка: %v", err)
			return &orderv1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "Что-то пошло не так",
			}, nil
		}
	}

	return converter.ModelToPayResponse(order), nil
}
