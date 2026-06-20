package grpc

import (
	"context"

	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
)

type InventoryClient interface {
	ListParts(ctx context.Context, in *inventoryv1.ListPartsRequest, opts ...grpc.CallOption) (*inventoryv1.ListPartsResponse, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, in *paymentv1.PayOrderRequest, opts ...grpc.CallOption) (*paymentv1.PayOrderResponse, error)
}
