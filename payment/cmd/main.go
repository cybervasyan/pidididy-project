package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	paymentv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type paymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (p *paymentService) PayOrder(context.Context, *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	return &paymentv1.PayOrderResponse{
		TransactionUuid: uuid.New().String(),
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50052))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()

	// Создаем gRPC сервер
	s := grpc.NewServer()

	// Регистрируем наш сервис
	service := &paymentService{}

	paymentv1.RegisterPaymentServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("PaymentService gRPC server listening on %d\n", 50052)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
