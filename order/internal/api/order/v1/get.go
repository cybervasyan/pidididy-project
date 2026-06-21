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

func (a *api) GetOrderByUuid(ctx context.Context, params orderv1.GetOrderByUuidParams) (r orderv1.GetOrderByUuidRes, _ error) {
	order, err := a.orderService.GetOrderByUuid(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderDoesntExist):
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Заказ не найден",
			}, nil
		default:
			log.Printf("GetOrderByUuid: непредвиденная ошибка для заказа %s: %v", params.OrderUUID, err)
			return &orderv1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "Что-то пошло не так",
			}, nil
		}
	}
	return converter.ModelToOrderDto(order), nil
}
