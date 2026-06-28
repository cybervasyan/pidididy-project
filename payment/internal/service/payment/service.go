package payment

import def "github.com/cybervasyan/pdididy-project/payment/internal/service"

var _ def.Payment = (*service)(nil)

type service struct {
}

func NewPaymentService() *service {
	return &service{}
}
