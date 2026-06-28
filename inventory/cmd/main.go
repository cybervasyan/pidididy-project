package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	inventoryV1 "github.com/cybervasyan/pdididy-project/inventory/internal/api/inventory/v1"
	repoPart "github.com/cybervasyan/pdididy-project/inventory/internal/repository/part"
	servPart "github.com/cybervasyan/pdididy-project/inventory/internal/service/part"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

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

	s := grpc.NewServer()

	repository := repoPart.NewRepository(seedParts())
	service := servPart.NewPartService(repository)
	invAPI := inventoryV1.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(s, invAPI)

	reflection.Register(s)

	go func() {
		log.Printf("InventoryService gRPC server listening on %d\n", 50051)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server stopped")
}
