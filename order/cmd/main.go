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
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
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

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]*orderv1.OrderDto
}

type OrderHandler struct {
	orderv1.UnimplementedHandler
	paymentClient   paymentv1.PaymentServiceClient
	inventoryClient inventoryv1.InventoryServiceClient
	storage         *OrderStorage
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[uuid.UUID]*orderv1.OrderDto),
	}
}

func NewOrderHandler(storage *OrderStorage, pay paymentv1.PaymentServiceClient, inv inventoryv1.InventoryServiceClient) *OrderHandler {
	return &OrderHandler{
		storage:         storage,
		paymentClient:   pay,
		inventoryClient: inv,
	}
}

func (o *OrderHandler) GetOrderByUuid(_ context.Context, params orderv1.GetOrderByUuidParams) (r orderv1.GetOrderByUuidRes, _ error) {
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

func (o *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (r orderv1.CreateOrderRes, _ error) {
	o.storage.mu.Lock()

	orderToCreate := orderv1.OrderDto{
		OrderUUID: uuid.New(),
		UserUUID:  req.UserUUID,
		PartUuids: req.PartUuids,
		Status:    orderv1.OrderStatusPENDINGPAYMENT,
	}

	o.storage.orders[orderToCreate.OrderUUID] = &orderToCreate
	order, ok := o.storage.orders[orderToCreate.OrderUUID]
	if !ok {
		o.storage.mu.Unlock()
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Такого элемента нет",
		}, nil
	}

	var partsStringUUIDs []string

	for _, part := range order.PartUuids {
		partsStringUUIDs = append(partsStringUUIDs, part.String())
	}

	o.storage.mu.Unlock()
	parts, err := o.inventoryClient.ListParts(ctx, &inventoryv1.ListPartsRequest{Filter: &inventoryv1.PartsFilter{Uuids: partsStringUUIDs}})
	if err != nil {
		return &orderv1.InternalServerError{
			Code:    500,
			Message: "InventoryService не отвечает",
		}, nil
	}
	o.storage.mu.Lock()
	var totalPrice float64

	for _, part := range parts.Parts {
		totalPrice += part.Price
	}

	order.SetTotalPrice(totalPrice)
	o.storage.mu.Unlock()
	return &orderv1.CreateOrderResponse{
		OrderUUID: orderToCreate.GetOrderUUID(),
	}, nil
}

func (o *OrderHandler) CancelOrder(_ context.Context, params orderv1.CancelOrderParams) (r orderv1.CancelOrderRes, _ error) {
	o.storage.mu.Lock()
	defer o.storage.mu.Unlock()

	order, ok := o.storage.orders[params.OrderUUID]
	if !ok {
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Такого элемента нет",
		}, nil
	}

	if order.GetStatus() != orderv1.OrderStatusPENDINGPAYMENT {
		return &orderv1.ConflictError{
			Code:    422,
			Message: "Заказ находится в статусе отличном от PENDING",
		}, nil
	}

	order.SetStatus(orderv1.OrderStatusCANCELLED)
	return &orderv1.CancelOrderNoContent{}, nil
}

func (o *OrderHandler) PayOrder(ctx context.Context, req *orderv1.PayOrderRequest, params orderv1.PayOrderParams) (r orderv1.PayOrderRes, _ error) {
	o.storage.mu.Lock()

	order, ok := o.storage.orders[params.OrderUUID]

	if !ok {
		o.storage.mu.Unlock()
		return &orderv1.NotFoundError{
			Code:    404,
			Message: "Такого элемента нет",
		}, nil
	}

	if order.GetStatus() != orderv1.OrderStatusPENDINGPAYMENT {
		o.storage.mu.Unlock()
		return &orderv1.ConflictError{
			Code:    422,
			Message: "Заказ находится в статусе отличном от PENDING",
		}, nil
	}
	var paymentMethod paymentv1.PaymentMethod

	switch req.PaymentMethod {
	case orderv1.PaymentMethodCARD:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderv1.PaymentMethodSBP:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderv1.PaymentMethodCREDITCARD:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderv1.PaymentMethodINVESTORMONEY:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case orderv1.PaymentMethodUNKNOWN:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	default:
		paymentMethod = paymentv1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
	o.storage.mu.Unlock()
	transactionUUID, err := o.paymentClient.PayOrder(ctx, &paymentv1.PayOrderRequest{
		OrderUuid:     params.OrderUUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return &orderv1.InternalServerError{
			Code:    500,
			Message: "PaymentService не отвечает",
		}, nil
	}
	parsedUUID, err := uuid.Parse(transactionUUID.GetTransactionUuid())
	if err != nil {
		return &orderv1.InternalServerError{
			Code:    500,
			Message: "PaymentService сломался",
		}, nil
	}

	o.storage.mu.Lock()
	order.SetStatus(orderv1.OrderStatusPAID)
	order.SetPaymentMethod(orderv1.NewOptPaymentMethod(req.PaymentMethod))
	order.SetTransactionUUID(orderv1.NewOptNilUUID(parsedUUID))
	o.storage.mu.Unlock()
	return &orderv1.PayOrderResponse{
		TransactionUUID: parsedUUID,
	}, nil
}

func main() {
	storage := NewOrderStorage()
	pay, err := grpc.NewClient(paymentServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка создания gRPC сервера PaymentService: %v", err)
	}

	inv, err := grpc.NewClient(inventoryServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Ошибка создания gRPC сервера InventoryService: %v", err)
	}

	defer func() {
		if cerr := pay.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	defer func() {
		if cerr := inv.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	payClient := paymentv1.NewPaymentServiceClient(pay)
	invClient := inventoryv1.NewInventoryServiceClient(inv)

	orderHandler := NewOrderHandler(storage, payClient, invClient)

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
