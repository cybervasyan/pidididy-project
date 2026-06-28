package v1

import (
	"context"
	"log"

	"github.com/cybervasyan/pdididy-project/inventory/internal/converter"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	filter, err := converter.PartsFilterToModel(req.GetFilter())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "невалидный фильтр: %v", err)
	}

	parts, err := a.inventoryService.ListParts(ctx, filter)
	if err != nil {
		log.Printf("ListParts: непредвиденная ошибка: %v", err)
		return nil, status.Error(codes.Internal, "что-то пошло не так")
	}

	return &inventoryv1.ListPartsResponse{
		Parts: converter.PartsToProto(parts),
	}, nil
}
