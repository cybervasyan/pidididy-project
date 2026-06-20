package converter

import (
	"github.com/cybervasyan/pdididy-project/order/internal/model"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
)

func PaymentMethodToProto(method *model.PaymentMethod) paymentv1.PaymentMethod {
	if method == nil {
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}

	switch *method {
	case model.PaymentMethodCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCREDITCARD:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
