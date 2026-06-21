package order

import (
	"context"

	"github.com/cybervasyan/pdididy-project/order/internal/model"
	repoConverter "github.com/cybervasyan/pdididy-project/order/internal/repository/converter"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
)

func (s *service) CreateOrder(ctx context.Context, req *model.Order) (model.Order, error) {
	partsStringUUIDs := make([]string, 0, len(req.PartUuids))
	for _, id := range req.PartUuids {
		partsStringUUIDs = append(partsStringUUIDs, id.String())
	}

	parts, err := s.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{Filter: &inventoryv1.PartsFilter{Uuids: partsStringUUIDs}})
	if err != nil {
		return model.Order{}, err
	}

	foundParts := make(map[string]struct{}, len(parts.Parts))
	for _, part := range parts.Parts {
		foundParts[part.Uuid] = struct{}{}
	}

	for _, id := range req.PartUuids {
		if _, ok := foundParts[id.String()]; !ok {
			return model.Order{}, model.ErrPartNotFound
		}
	}
	var totalPrice float64

	for _, part := range parts.Parts {
		totalPrice += part.Price
	}

	newOrder := *req
	newOrder.OrderUUID = uuid.New()
	newOrder.TotalPrice = totalPrice
	newOrder.Status = model.OrderStatusPENDINGPAYMENT

	repoModel := repoConverter.ToRepoModel(newOrder)

	err = s.orderRepo.Create(ctx, &repoModel)
	if err != nil {
		return model.Order{}, err
	}

	return repoConverter.ToServiceModel(repoModel), nil
}
