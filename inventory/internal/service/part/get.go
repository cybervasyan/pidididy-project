package part

import (
	"context"
	"errors"

	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/converter"
	repoModel "github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
)

func (s *service) GetPart(ctx context.Context, partUUID uuid.UUID) (model.Part, error) {
	part, err := s.partRepo.Get(ctx, partUUID)

	if err != nil {
		if errors.Is(err, repoModel.ErrPartNotFound) {
			return model.Part{}, model.ErrPartNotFound
		}

		return model.Part{}, err
	}

	return converter.PartToServiceModel(part), nil
}
