package model

import (
	"time"

	"github.com/google/uuid"
)

type PartsFilter struct {
	PartUUIDs             []uuid.UUID
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

type Part struct {
	PartUUID      uuid.UUID
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]Value
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Category string

const (
	CategoryUnspecified Category = "CATEGORY_UNSPECIFIED"
	CategoryEngine      Category = "CATEGORY_ENGINE"
	CategoryFuel        Category = "CATEGORY_FUEL"
	CategoryPorthole    Category = "CATEGORY_PORTHOLE"
	CategoryWing        Category = "CATEGORY_WING"
)

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

// ValueKind говорит, какое из полей Value валидно (эмуляция proto oneof).
type ValueKind string

const (
	ValueKindUnspecified ValueKind = ""
	ValueKindString      ValueKind = "STRING"
	ValueKindInt64       ValueKind = "INT64"
	ValueKindDouble      ValueKind = "DOUBLE"
	ValueKindBool        ValueKind = "BOOL"
)

type Value struct {
	Kind        ValueKind
	StringValue string
	Int64Value  int64
	DoubleValue float64
	BoolValue   bool
}
