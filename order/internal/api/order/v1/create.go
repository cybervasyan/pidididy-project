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

func (a *api) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (r orderv1.CreateOrderRes, _ error) {
	orderModel := converter.CreateRequestToModel(req)

	created, err := a.orderService.CreateOrder(ctx, &orderModel)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartNotFound):
			return &orderv1.NotFoundError{
				Code:    http.StatusNotFound,
				Message: "Одна или несколько запчастей не найдены",
			}, nil

		default:
			log.Printf("CreateOrder: непредвиденная ошибка: %v", err)
			return &orderv1.InternalServerError{
				Code:    http.StatusInternalServerError,
				Message: "Что-то пошло не так",
			}, nil
		}
	}

	return converter.ModelToCreateResponse(created), nil
}
