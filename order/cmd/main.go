package main

import (
	"context"
	"log"
	"sync"
	"time"

	orderv1 "github.com/cybervasyan/pdididy-project/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

const (
	httpPort  = "8080"
	urlParams = "city"
	// Таймауты для http-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderv1.OrderDto
}

type OrderHandler struct {
	storage *OrderStorage
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderv1.OrderDto),
	}
}

func NewOrderHandler(storage *OrderStorage) *OrderHandler {
	return &OrderHandler{
		storage: storage,
	}
}

func (o *OrderStorage) GetOrderByUuid(ctx context.Context, params orderv1.GetOrderByUuidParams) (r orderv1.GetOrderByUuidRes, _ error) {
	o.mu.RLock()
	defer o.mu.RUnlock()

	order, ok := o.orders[params.OrderUUID]
	if !ok {
		return &orderv1.NotFoundError{}, nil
	}

	return order, nil
}

func main() {
	storage := NewOrderStorage()

	orderHandler := NewOrderHandler(storage)

	orderServer, err := orderv1.NewServer(orderHandler)
	if err != nil {
		log.Fatalf("Ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()
}
