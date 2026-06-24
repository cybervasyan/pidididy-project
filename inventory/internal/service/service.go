package service

import (
	"context"

	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	"github.com/google/uuid"
)

type Part interface {
	GetPart(ctx context.Context, partUUID uuid.UUID) (model.Part, error)
	ListParts(ctx context.Context, req model.PartsFilter) ([]model.Part, error)
}
