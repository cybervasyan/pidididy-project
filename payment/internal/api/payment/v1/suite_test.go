package v1

import (
	"context"
	"testing"

	"github.com/cybervasyan/pdididy-project/payment/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type APISuite struct {
	suite.Suite

	ctx context.Context

	paymentService *mocks.MockPayment

	api *api
}

func (s *APISuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentService = mocks.NewMockPayment(s.T())

	s.api = NewAPI(
		s.paymentService,
	)
}

func TestAPI(t *testing.T) {
	suite.Run(t, new(APISuite))
}
