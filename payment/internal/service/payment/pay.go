package payment

import (
	"context"

	"github.com/cybervasyan/pdididy-project/payment/internal/model"
	"github.com/google/uuid"
)

func (s *service) PayOrder(_ context.Context, orderUUID uuid.UUID, userUUID uuid.UUID, paymentMethod model.PaymentMethod) (uuid.UUID, error) {
	if orderUUID == uuid.Nil {
		return uuid.Nil, model.ErrOrderUUIDRequired
	}

	if userUUID == uuid.Nil {
		return uuid.Nil, model.ErrUserUUIDRequired
	}

	if paymentMethod == model.PaymentMethodUnspecified {
		return uuid.Nil, model.ErrPaymentMethodUnspecified
	}

	return uuid.New(), nil
}
