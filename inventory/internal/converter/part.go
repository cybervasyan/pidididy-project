package converter

import (
	"fmt"

	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	inventoryv1 "github.com/cybervasyan/pdididy-project/shared/pkg/proto/inventory/v1"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PartsFilterToModel(f *inventoryv1.PartsFilter) (model.PartsFilter, error) {
	if f == nil {
		return model.PartsFilter{}, nil
	}

	uuids := make([]uuid.UUID, 0, len(f.GetUuids()))
	for _, raw := range f.GetUuids() {
		id, err := uuid.Parse(raw)
		if err != nil {
			return model.PartsFilter{}, fmt.Errorf("uuid %q: %w", raw, err)
		}

		uuids = append(uuids, id)
	}

	categories := make([]model.Category, 0, len(f.GetCategories()))
	for _, c := range f.GetCategories() {
		categories = append(categories, categoryToModel(c))
	}

	return model.PartsFilter{
		PartUUIDs:             uuids,
		Names:                 f.GetNames(),
		Categories:            categories,
		ManufacturerCountries: f.GetManufacturerCountries(),
		Tags:                  f.GetTags(),
	}, nil
}

func categoryToModel(c inventoryv1.Category) model.Category {
	switch c {
	case inventoryv1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryv1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryv1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryv1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func PartToProto(p model.Part) *inventoryv1.Part {
	metadata := make(map[string]*inventoryv1.Value, len(p.Metadata))
	for k, v := range p.Metadata {
		metadata[k] = valueToProto(v)
	}

	return &inventoryv1.Part{
		Uuid:          p.PartUUID.String(),
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      categoryToProto(p.Category),
		Dimensions:    dimensionsToProto(p.Dimensions),
		Manufacturer:  manufacturerToProto(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      metadata,
		CreatedAt:     timestamppb.New(p.CreatedAt),
		UpdatedAt:     timestamppb.New(p.UpdatedAt),
	}
}

func PartsToProto(parts []model.Part) []*inventoryv1.Part {
	result := make([]*inventoryv1.Part, 0, len(parts))
	for _, p := range parts {
		result = append(result, PartToProto(p))
	}

	return result
}

func categoryToProto(c model.Category) inventoryv1.Category {
	switch c {
	case model.CategoryEngine:
		return inventoryv1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryv1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryv1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryv1.Category_CATEGORY_WING
	default:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
}

func dimensionsToProto(d *model.Dimensions) *inventoryv1.Dimensions {
	if d == nil {
		return nil
	}

	return &inventoryv1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func manufacturerToProto(m *model.Manufacturer) *inventoryv1.Manufacturer {
	if m == nil {
		return nil
	}

	return &inventoryv1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func valueToProto(v model.Value) *inventoryv1.Value {
	switch v.Kind {
	case model.ValueKindString:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_StringValue{StringValue: v.StringValue}}
	case model.ValueKindInt64:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_Int64Value{Int64Value: v.Int64Value}}
	case model.ValueKindDouble:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_DoubleValue{DoubleValue: v.DoubleValue}}
	case model.ValueKindBool:
		return &inventoryv1.Value{Kind: &inventoryv1.Value_BoolValue{BoolValue: v.BoolValue}}
	default:
		return &inventoryv1.Value{}
	}
}
