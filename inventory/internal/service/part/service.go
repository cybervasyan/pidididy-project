package part

import "github.com/cybervasyan/pdididy-project/inventory/internal/repository"

type service struct {
	partRepo repository.Repository
}

func NewPartService(partRepo repository.Repository) *service {
	return &service{
		partRepo: partRepo,
	}
}
