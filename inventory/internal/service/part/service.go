package part

import (
	"github.com/cybervasyan/pdididy-project/inventory/internal/repository"
	def "github.com/cybervasyan/pdididy-project/inventory/internal/service"
)

var _ def.Part = (*service)(nil)

type service struct {
	partRepo repository.Repository
}

func NewPartService(partRepo repository.Repository) *service {
	return &service{
		partRepo: partRepo,
	}
}
