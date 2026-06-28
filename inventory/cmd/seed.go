package main

import (
	"time"

	repoModel "github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
)

func seedParts() []repoModel.Part {
	now := time.Now()

	return []repoModel.Part{
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
			Name:          "Main Engine",
			Description:   "Основной маршевый двигатель",
			Price:         100.0,
			StockQuantity: 5,
			Category:      repoModel.CategoryEngine,
			Dimensions: &repoModel.Dimensions{
				Length: 4.2,
				Width:  1.8,
				Height: 1.8,
				Weight: 1200.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Rocketdyne",
				Country: "USA",
				Website: "https://rocketdyne.example",
			},
			Tags: []string{"engine", "main"},
			Metadata: map[string]repoModel.Value{
				"thrust_kn": {Kind: repoModel.ValueKindDouble, DoubleValue: 845.0},
				"reusable":  {Kind: repoModel.ValueKindBool, BoolValue: true},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
			Name:          "Fuel Tank",
			Description:   "Бак для топлива",
			Price:         50.0,
			StockQuantity: 12,
			Category:      repoModel.CategoryFuel,
			Dimensions: &repoModel.Dimensions{
				Length: 6.0,
				Width:  3.0,
				Height: 3.0,
				Weight: 800.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Orbital Systems",
				Country: "Germany",
				Website: "https://orbital.example",
			},
			Tags: []string{"fuel", "tank"},
			Metadata: map[string]repoModel.Value{
				"capacity_l": {Kind: repoModel.ValueKindInt64, Int64Value: 25000},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			PartUUID:      uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
			Name:          "Porthole",
			Description:   "Иллюминатор обзорный",
			Price:         15.0,
			StockQuantity: 30,
			Category:      repoModel.CategoryPorthole,
			Tags:          []string{"porthole", "glass"},
			Metadata: map[string]repoModel.Value{
				"material": {Kind: repoModel.ValueKindString, StringValue: "quartz"},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}
