package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	orderV1 "github.com/cybervasyan/pdididy-project/order/internal/api/order/v1"
	customMiddleware "github.com/cybervasyan/pdididy-project/order/internal/middleware"
	repoOrder "github.com/cybervasyan/pdididy-project/order/internal/repository/order"
	servOrder "github.com/cybervasyan/pdididy-project/order/internal/service/order"
	orderv1 "github.com/cybervasyan/pdididy-project/shared/pkg/openapi/order/v1"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	httpPort = "8080"
	// Таймауты для http-сервера
	readHeaderTimeout      = 5 * time.Second
	shutdownTimeout        = 10 * time.Second
	paymentServerAddress   = "localhost:50052"
	inventoryServerAddress = "localhost:50051"
)

func main() {
	pay, err := grpc.NewClient(paymentServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Ошибка создания gRPC сервера PaymentService: %v", err)
		return
	}

	defer func() {
		if gerr := pay.Close(); gerr != nil {
			log.Printf("failed to close connect: %v", gerr)
		}
	}()

	inv, err := grpc.NewClient(inventoryServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Ошибка создания gRPC сервера InventoryService: %v", err)
		return
	}

	defer func() {
		if gerr := inv.Close(); gerr != nil {
			log.Printf("failed to close connect: %v", gerr)
		}
	}()

	payClient := paymentv1.NewPaymentServiceClient(pay)
	invClient := inventoryv1.NewInventoryServiceClient(inv)

	repository := repoOrder.NewRepository()
	service := servOrder.NewOrderService(repository, invClient, payClient)

	orderHandler := orderV1.NewAPI(service)

	orderServer, err := orderv1.NewServer(orderHandler)
	if err != nil {
		log.Printf("Ошибка создания сервера OpenAPI: %v", err)
		return
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
		gerr := server.ListenAndServe()
		if gerr != nil && !errors.Is(gerr, http.ErrServerClosed) {
			log.Printf("Ошибка запуска сервера: %v\n", gerr)
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
