package payment

import (
	"context"
	"testing"

	"github.com/cybervasyan/pdididy-project/payment/internal/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestPayOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		orderUUID     uuid.UUID
		userUUID      uuid.UUID
		paymentMethod model.PaymentMethod
		wantErr       error
	}{
		{
			name:          "success",
			orderUUID:     uuid.New(),
			userUUID:      uuid.New(),
			paymentMethod: model.PaymentMethodCard,
			wantErr:       nil,
		},
		{
			name:          "empty order uuid",
			orderUUID:     uuid.Nil,
			userUUID:      uuid.New(),
			paymentMethod: model.PaymentMethodCard,
			wantErr:       model.ErrOrderUUIDRequired,
		},
		{
			name:          "empty user uuid",
			orderUUID:     uuid.New(),
			userUUID:      uuid.Nil,
			paymentMethod: model.PaymentMethodCard,
			wantErr:       model.ErrUserUUIDRequired,
		},
		{
			name:          "unspecified payment method",
			orderUUID:     uuid.New(),
			userUUID:      uuid.New(),
			paymentMethod: model.PaymentMethodUnspecified,
			wantErr:       model.ErrPaymentMethodUnspecified,
		},
	}

	svc := NewPaymentService()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			transactionUUID, err := svc.PayOrder(context.Background(), tt.orderUUID, tt.userUUID, tt.paymentMethod)

			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, uuid.Nil, transactionUUID)
				return
			}

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil, transactionUUID)
		})
	}
}
