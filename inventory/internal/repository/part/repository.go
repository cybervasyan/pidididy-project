package part

import (
	"sync"

	def "github.com/cybervasyan/pdididy-project/inventory/internal/repository"
	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
)

var _ def.Repository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[uuid.UUID]model.Part
}

func NewRepository(parts []model.Part) *repository {
	m := make(map[uuid.UUID]model.Part, len(parts))
	for _, p := range parts {
		m[p.PartUUID] = p
	}

	return &repository{
		parts: m,
	}
}
