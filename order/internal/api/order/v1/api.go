package v1

import (
	"github.com/cybervasyan/pdididy-project/order/internal/service"
	orderv1 "github.com/cybervasyan/pdididy-project/shared/pkg/openapi/order/v1"
)

type api struct {
	orderv1.UnimplementedHandler

	orderService service.Order
}

func NewAPI(orderService service.Order) *api {
	return &api{
		orderService: orderService,
	}
}
