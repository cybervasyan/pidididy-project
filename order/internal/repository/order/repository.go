package order

import (
	"sync"

	"github.com/cybervasyan/pdididy-project/order/internal/repository/model"
	"github.com/google/uuid"
)

type repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*model.Order
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[uuid.UUID]*model.Order),
	}
}
