package part

import (
	"context"

	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/converter"
)

func (s *service) ListParts(ctx context.Context, req model.PartsFilter) ([]model.Part, error) {
	parts, err := s.partRepo.List(ctx, converter.PartsFilterToRepoModel(req))

	if err != nil {
		return []model.Part{}, err
	}

	return converter.PartsToServiceModel(parts), nil
}
