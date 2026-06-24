package part

import (
	"context"

	"github.com/cybervasyan/pdididy-project/inventory/internal/repository/model"
	"github.com/google/uuid"
)

func (r *repository) List(_ context.Context, req model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []model.Part

	for _, v := range r.parts {
		parts = append(parts, v)
	}

	if len(req.PartUUIDs) != 0 {
		parts = filterPartsByUUIDs(parts, req.PartUUIDs)
	}

	if len(req.Names) != 0 {
		parts = filterPartsByNames(parts, req.Names)
	}

	if len(req.Categories) != 0 {
		parts = filterPartsByCategories(parts, req.Categories)
	}

	if len(req.ManufacturerCountries) != 0 {
		parts = filterPartsByCountry(parts, req.ManufacturerCountries)
	}

	if len(req.Tags) != 0 {
		parts = filterPartsByTags(parts, req.Tags)
	}

	return parts, nil
}

func filterPartsByUUIDs(parts []model.Part, uuids []uuid.UUID) []model.Part {
	set := make(map[uuid.UUID]struct{}, len(uuids))
	for _, v := range uuids {
		set[v] = struct{}{}
	}
	var result []model.Part
	for _, p := range parts {
		if _, ok := set[p.PartUUID]; ok {
			result = append(result, p)
		}
	}
	return result
}

func filterPartsByNames(parts []model.Part, names []string) []model.Part {
	set := make(map[string]struct{}, len(names))
	for _, v := range names {
		set[v] = struct{}{}
	}
	var result []model.Part
	for _, p := range parts {
		if _, ok := set[p.Name]; ok {
			result = append(result, p)
		}
	}
	return result
}

func filterPartsByCategories(parts []model.Part, cats []model.Category) []model.Part {
	set := make(map[model.Category]struct{}, len(cats))
	for _, v := range cats {
		set[v] = struct{}{}
	}
	var result []model.Part
	for _, p := range parts {
		if _, ok := set[p.Category]; ok {
			result = append(result, p)
		}
	}
	return result
}

func filterPartsByCountry(parts []model.Part, countries []string) []model.Part {
	set := make(map[string]struct{}, len(countries))
	for _, v := range countries {
		set[v] = struct{}{}
	}
	var result []model.Part
	for _, p := range parts {
		if p.Manufacturer != nil {
			if _, ok := set[p.Manufacturer.Country]; ok {
				result = append(result, p)
			}
		}
	}
	return result
}

func filterPartsByTags(parts []model.Part, tags []string) []model.Part {
	set := make(map[string]struct{}, len(tags))
	for _, v := range tags {
		set[v] = struct{}{}
	}
	var result []model.Part
	for _, p := range parts {
		for _, tag := range p.Tags {
			if _, ok := set[tag]; ok {
				result = append(result, p)
				break
			}
		}
	}
	return result
}
