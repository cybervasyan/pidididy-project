package v1

import (
	"context"
	"log"
	"net/http"

	"github.com/cybervasyan/pdididy-project/order/internal/model"
	orderv1 "github.com/cybervasyan/pdididy-project/shared/pkg/openapi/order/v1"
	"github.com/go-faster/errors"
)

func (a *api) CancelOrder(ctx context.Context, params orderv1.CancelOrderParams) (r orderv1.CancelOrderRes, _ error) {
	err := a.orderService.CancelOrder(ctx, params.OrderUUID)
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
				Message: "Заказ находится в статусе, отличном от PENDING",
			}, nil

		default:
			log.Printf("CancelOrder: непредвиденная ошибка для заказа %s: %v", params.OrderUUID, err)
			return &orderv1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "Что-то пошло не так",
			}, nil
		}
	}

	return &orderv1.CancelOrderNoContent{}, nil
}
