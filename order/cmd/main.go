package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	customMiddleware "github.com/cybervasyan/pdididy-project/order/internal/middleware"
	orderv1 "github.com/cybervasyan/pdididy-project/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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
	orderv1.UnimplementedHandler
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

func (o *OrderHandler) GetOrderByUuid(ctx context.Context, params orderv1.GetOrderByUuidParams) (r orderv1.GetOrderByUuidRes, _ error) {
	o.storage.mu.RLock()
	defer o.storage.mu.RUnlock()

	order, ok := o.storage.orders[params.OrderUUID]
	if !ok {
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Такого элемента нет",
		}, nil
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

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(10 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              net.JoinHostPort("localhost", httpPort),
		Handler:           r,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("http-сервер запущен на порту %s\n", httpPort)
		err = server.ListenAndServe()
		if err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Завершение работы сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
	}

	log.Printf("Сервер остановлен")
}
