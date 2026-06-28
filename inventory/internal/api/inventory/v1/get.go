package v1

import (
	"context"
	"errors"
	"log"

	"github.com/cybervasyan/pdididy-project/inventory/internal/converter"
	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	partUUID, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "невалидный uuid: %v", err)
	}

	part, err := a.inventoryService.GetPart(ctx, partUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartNotFound):
			return nil, status.Error(codes.NotFound, "деталь не найдена")
		default:
			log.Printf("GetPart: непредвиденная ошибка для детали %s: %v", req.GetUuid(), err)
			return nil, status.Error(codes.Internal, "что-то пошло не так")
		}
	}

	return &inventoryv1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
