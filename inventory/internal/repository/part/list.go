package part

import (
	"context"

	"inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

// ListParts возвращает отфильтрованный список деталей.
// Использует pipeline для последовательной фильтрации с AND-логикой между полями
// и OR-логикой внутри каждого поля.
func (r *repository) ListParts(_ context.Context, filter *model.PartsFilter) []*model.Part {
	parts := r.Values()

	if filterIsEmpty(filter) {
		return parts
	}

	repoParts := converter.PartsToRepoModel(parts)
	repoFilter := converter.PartsFilterToRepoModel(*filter)

	// Строим и применяем pipeline фильтров
	pipeline := buildFilterPipeline(&repoFilter)
	filteredParts := applyFilterPipeline(repoParts, pipeline)

	return converter.PartsToModel(filteredParts)
}

// filterIsEmpty проверяет, пустой ли фильтр.
func filterIsEmpty(f *model.PartsFilter) bool {
	if f == nil {
		return true
	}
	return len(f.Uuids) == 0 &&
		len(f.Names) == 0 &&
		len(f.Categories) == 0 &&
		len(f.ManufacturerCountries) == 0 &&
		len(f.Tags) == 0
}

// FilterFunc представляет одну стадию фильтрации в pipeline.
type FilterFunc func([]*repoModel.Part) []*repoModel.Part

// applyFilterPipeline последовательно применяет все фильтры.
func applyFilterPipeline(parts []*repoModel.Part, pipeline []FilterFunc) []*repoModel.Part {
	for _, filter := range pipeline {
		parts = filter(parts)
		if len(parts) == 0 {
			return parts
		}
	}
	return parts
}

// buildFilterPipeline создаёт цепочку фильтров на основе PartsFilter.
func buildFilterPipeline(f *repoModel.PartsFilter) []FilterFunc {
	var pipeline []FilterFunc

	if len(f.Uuids) > 0 {
		pipeline = append(pipeline, func(parts []*repoModel.Part) []*repoModel.Part {
			return filterByField(parts, f.Uuids, func(p *repoModel.Part) string { return p.Uuid })
		})
	}

	if len(f.Names) > 0 {
		pipeline = append(pipeline, func(parts []*repoModel.Part) []*repoModel.Part {
			return filterByField(parts, f.Names, func(p *repoModel.Part) string { return p.Name })
		})
	}

	if len(f.Categories) > 0 {
		pipeline = append(pipeline, func(parts []*repoModel.Part) []*repoModel.Part {
			return filterByField(parts, f.Categories, func(p *repoModel.Part) repoModel.Category { return p.Category })
		})
	}

	if len(f.ManufacturerCountries) > 0 {
		pipeline = append(pipeline, func(parts []*repoModel.Part) []*repoModel.Part {
			return filterByField(parts, f.ManufacturerCountries, func(p *repoModel.Part) string {
				if p.Manufacturer == nil {
					return ""
				}
				return p.Manufacturer.Country
			})
		})
	}

	if len(f.Tags) > 0 {
		pipeline = append(pipeline, func(parts []*repoModel.Part) []*repoModel.Part {
			return filterByTags(parts, f.Tags)
		})
	}

	return pipeline
}

// filterByField возвращает детали, у которых значение поля есть в values.
//
// В отличие от реализации через slices.Contains (O(n²)),
// использует внутренний set на основе map для поиска за O(1),
// что обеспечивает общую сложность O(n + m).
//
// n — количество деталей, m — количество элементов фильтра.
func filterByField[T comparable](
	parts []*repoModel.Part,
	values []T,
	getField func(*repoModel.Part) T,
) []*repoModel.Part {
	if len(values) == 0 {
		return parts
	}

	set := toSet(values)
	result := make([]*repoModel.Part, 0, len(parts))

	for _, p := range parts {
		if _, exists := set[getField(p)]; exists {
			result = append(result, p)
		}
	}

	return result
}

// filterByTags оставляет детали с хотя бы одним тегом из списка (OR-логика).
func filterByTags(parts []*repoModel.Part, tags []string) []*repoModel.Part {
	if len(tags) == 0 {
		return parts
	}

	tagSet := toSet(tags)
	result := make([]*repoModel.Part, 0, len(parts))

	for _, p := range parts {
		if hasAnyTag(p.Tags, tagSet) {
			result = append(result, p)
		}
	}

	return result
}

// hasAnyTag проверяет наличие хотя бы одного тега из set'а.
func hasAnyTag(partTags []string, tagSet map[string]struct{}) bool {
	for _, tag := range partTags {
		if _, exists := tagSet[tag]; exists {
			return true
		}
	}
	return false
}

// toSet преобразует slice в set для O(1) поиска.
func toSet[T comparable](values []T) map[T]struct{} {
	set := make(map[T]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}
	return set
}
