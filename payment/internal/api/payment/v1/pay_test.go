package v1

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/cybervasyan/pdididy-project/payment/internal/model"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *APISuite) TestPaySuccess() {
	var (
		orderUUID       = uuid.New()
		userUUID        = uuid.New()
		paymentMethod   = model.PaymentMethodCard
		transactionUUID = uuid.New()
	)

	s.paymentService.EXPECT().PayOrder(s.ctx, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil)
	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().NoError(err)
	s.Require().NotNil(res)
	s.Require().Equal(transactionUUID.String(), res.TransactionUuid)
}

func (s *APISuite) TestPayWrongOrderUuid() {
	var (
		orderUUID = "suck on this"
		userUUID  = uuid.New()
	)

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
}

func (s *APISuite) TestPayWrongUserUuid() {
	var (
		orderUUID = uuid.New()
		userUUID  = "suck on this"
	)

	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID,
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
}

func (s *APISuite) TestPayUserUuidNotFound() {
	var (
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		paymentMethod = model.PaymentMethodCard
	)

	s.paymentService.EXPECT().PayOrder(s.ctx, orderUUID, userUUID, paymentMethod).Return(uuid.Nil, model.ErrUserUUIDRequired)
	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.InvalidArgument, st.Code())
}

func (s *APISuite) TestPayRandomError() {
	var (
		orderUUID     = uuid.New()
		userUUID      = uuid.New()
		paymentMethod = model.PaymentMethodCard
	)

	s.paymentService.EXPECT().PayOrder(s.ctx, orderUUID, userUUID, paymentMethod).Return(uuid.Nil, gofakeit.Error())
	req := &paymentv1.PayOrderRequest{
		OrderUuid:     orderUUID.String(),
		UserUuid:      userUUID.String(),
		PaymentMethod: paymentv1.PaymentMethod_PAYMENT_METHOD_CARD,
	}

	res, err := s.api.PayOrder(s.ctx, req)
	s.Require().Error(err)
	s.Require().Nil(res)

	st, ok := status.FromError(err)
	s.Require().True(ok)
	s.Require().Equal(codes.Internal, st.Code())
}
