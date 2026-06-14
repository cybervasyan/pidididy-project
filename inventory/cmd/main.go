package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer
	mu    sync.RWMutex
	parts []*inventoryv1.Part
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

func (s *inventoryService) ListParts(_ context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	filters := req.GetFilter()
	var parts []*inventoryv1.Part

	parts = append(parts, s.parts...)

	if filters == nil {
		return &inventoryv1.ListPartsResponse{
			Parts: parts,
		}, nil
	}

	if len(filters.GetUuids()) != 0 {
		parts = filterPartsByUUIDs(parts, filters.GetUuids())
	}

	if len(filters.GetNames()) != 0 {
		parts = filterPartsByNames(parts, filters.GetNames())
	}

	if len(filters.GetCategories()) != 0 {
		parts = filterPartsByCategories(parts, filters.GetCategories())
	}

	if len(filters.GetManufacturerCountries()) != 0 {
		parts = filterPartsByCountry(parts, filters.GetManufacturerCountries())
	}

	if len(filters.GetTags()) != 0 {
		parts = filterPartsByTags(parts, filters.GetTags())
	}

	return &inventoryv1.ListPartsResponse{
		Parts: parts,
	}, nil
}

func filterPartsByUUIDs(parts []*inventoryv1.Part, uuids []string) []*inventoryv1.Part {
	set := make(map[string]struct{}, len(uuids))
	for _, v := range uuids {
		set[v] = struct{}{}
	}
	var result []*inventoryv1.Part
	for _, p := range parts {
		if _, ok := set[p.GetUuid()]; ok {
			result = append(result, p)
		}
	}
	return result
}

func filterPartsByNames(parts []*inventoryv1.Part, names []string) []*inventoryv1.Part {
	set := make(map[string]struct{}, len(names))
	for _, v := range names {
		set[v] = struct{}{}
	}
	var result []*inventoryv1.Part
	for _, p := range parts {
		if _, ok := set[p.GetName()]; ok {
			result = append(result, p)
		}
	}
	return result
}

func filterPartsByCategories(parts []*inventoryv1.Part, cats []inventoryv1.Category) []*inventoryv1.Part {
	set := make(map[inventoryv1.Category]struct{}, len(cats))
	for _, v := range cats {
		set[v] = struct{}{}
	}
	var result []*inventoryv1.Part
	for _, p := range parts {
		if _, ok := set[p.GetCategory()]; ok {
			result = append(result, p)
		}
	}
	return result
}

func filterPartsByCountry(parts []*inventoryv1.Part, countries []string) []*inventoryv1.Part {
	set := make(map[string]struct{}, len(countries))
	for _, v := range countries {
		set[v] = struct{}{}
	}
	var result []*inventoryv1.Part
	for _, p := range parts {
		if p.GetManufacturer() != nil {
			if _, ok := set[p.GetManufacturer().GetCountry()]; ok {
				result = append(result, p)
			}
		}
	}
	return result
}

func filterPartsByTags(parts []*inventoryv1.Part, tags []string) []*inventoryv1.Part {
	set := make(map[string]struct{}, len(tags))
	for _, v := range tags {
		set[v] = struct{}{}
	}
	var result []*inventoryv1.Part
	for _, p := range parts {
		for _, tag := range p.GetTags() {
			if _, ok := set[tag]; ok {
				result = append(result, p)
				break
			}
		}
	}
	return result
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
	service := &inventoryService{
		parts: []*inventoryv1.Part{
			{
				Uuid:  "550e8400-e29b-41d4-a716-446655440000",
				Name:  "Main Engine",
				Price: 100.0,
			},
			{
				Uuid:  "550e8400-e29b-41d4-a716-446655440001",
				Name:  "Fuel Tank",
				Price: 50.0,
			},
		},
	}

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
