package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
}

func (i *inventoryService) GetPart(_ context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	return &inventoryv1.GetPartResponse{
		Part: &inventoryv1.Part{
			Uuid:  req.GetUuid(),
			Name:  "Test Part",
			Price: 100.0,
		},
	}, nil
}

func (i *inventoryService) ListParts(ctx context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	var parts []*inventoryv1.Part
	if req.Filter == nil || len(req.Filter.Uuids) == 0 {
		parts = append(parts, &inventoryv1.Part{Uuid: uuid.New().String(), Name: "Test Part", Price: 100.0})
		return &inventoryv1.ListPartsResponse{
			Parts: parts,
		}, nil
	}
	for index := range req.Filter.Uuids {
		parts = append(parts, &inventoryv1.Part{
			Price: 100.00,
			Name:  "Part N" + strconv.Itoa(index),
			Uuid:  req.Filter.Uuids[index],
		})
	}
	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
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
	service := &inventoryService{}

	inventoryv1.RegisterInventoryServiceServer(s, service)

	// Включаем рефлексию для отладки
	reflection.Register(s)

	go func() {
		log.Printf("InventoryService gRPC server listening on %d\n", 50051)
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
