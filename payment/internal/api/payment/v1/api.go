package v1

import (
	"github.com/cybervasyan/pdididy-project/payment/internal/service"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
)

var _ paymentv1.PaymentServiceServer = (*api)(nil)

type api struct {
	paymentv1.UnimplementedPaymentServiceServer

	paymentService service.Payment
}

func NewAPI(paymentService service.Payment) *api {
	return &api{
		paymentService: paymentService,
	}
}
