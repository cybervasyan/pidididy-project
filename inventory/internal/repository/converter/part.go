package converter

import (
	"github.com/cybervasyan/pdididy-project/inventory/internal/model"
	repoModel "github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
)

func convertValueToRepo(v model.Value) repoModel.Value {
	result := repoModel.Value{
		Kind:        repoModel.ValueKind(v.Kind),
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}

	return result
}

func convertValueToService(v repoModel.Value) model.Value {
	result := model.Value{
		Kind:        model.ValueKind(v.Kind),
		StringValue: v.StringValue,
		Int64Value:  v.Int64Value,
		DoubleValue: v.DoubleValue,
		BoolValue:   v.BoolValue,
	}

	return result
}

func convertDimensionToRepo(dim *model.Dimensions) *repoModel.Dimensions {
	if dim == nil {
		return nil
	}

	result := repoModel.Dimensions{
		Length: dim.Length,
		Weight: dim.Weight,
		Height: dim.Height,
		Width:  dim.Width,
	}

	return &result
}

func convertDimensionToService(dim *repoModel.Dimensions) *model.Dimensions {
	if dim == nil {
		return nil
	}

	result := model.Dimensions{
		Length: dim.Length,
		Weight: dim.Weight,
		Height: dim.Height,
		Width:  dim.Width,
	}

	return &result
}

func convertManufacturerToRepo(m *model.Manufacturer) *repoModel.Manufacturer {
	if m == nil {
		return nil
	}

	result := repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}

	return &result
}

func convertManufacturerToService(m *repoModel.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}

	result := model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}

	return &result
}

func PartToServiceModel(part repoModel.Part) model.Part {
	m := make(map[string]model.Value)

	for k, v := range part.Metadata {
		m[k] = convertValueToService(v)
	}

	result := model.Part{
		PartUUID:      part.PartUUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    convertDimensionToService(part.Dimensions),
		Manufacturer:  convertManufacturerToService(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      m,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}

	return result
}

func PartsToServiceModel(parts []repoModel.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))

	for _, p := range parts {
		result = append(result, PartToServiceModel(p))
	}

	return result
}

func PartsFilterToRepoModel(pf model.PartsFilter) repoModel.PartsFilter {
	c := make([]repoModel.Category, 0, len(pf.Categories))

	for _, v := range pf.Categories {
		c = append(c, repoModel.Category(v))
	}

	result := repoModel.PartsFilter{
		PartUUIDs:             pf.PartUUIDs,
		Names:                 pf.Names,
		Categories:            c,
		ManufacturerCountries: pf.ManufacturerCountries,
		Tags:                  pf.Tags,
	}

	return result
}
