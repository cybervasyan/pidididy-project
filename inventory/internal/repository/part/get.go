package part

import (
	"context"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) Get(_ context.Context, partUUID uuid.UUID) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[partUUID]

	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return part, nil
}
