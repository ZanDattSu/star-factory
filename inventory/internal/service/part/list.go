package part

import (
	"context"
	"errors"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
)

// ListParts возвращает отфильтрованный список деталей.
// Использует pipeline для последовательной фильтрации с AND-логикой между полями
// и OR-логикой внутри каждого поля.
func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	parts, err := s.repository.ListParts(ctx)
	if err != nil {
		return []*model.Part{}, errors.New("error listing parts")
	}

	if filterIsEmpty(filter) {
		return parts, nil
	}

	// Предварительно создаём set'ы для O(1) lookup
	uuidSet := toSet(filter.Uuids)
	nameSet := toSet(filter.Names)
	countrySet := toSet(filter.ManufacturerCountries)
	tagSet := toSet(filter.Tags)
	categorySet := toSet(filter.Categories)

	// Один проход по всем деталям с ранним выходом
	filteredParts := make([]*model.Part, 0, len(parts))
	for _, part := range parts {
		if matchesPart(
			part, filter, uuidSet, nameSet, countrySet, tagSet, categorySet,
		) {
			filteredParts = append(filteredParts, part)
		}
	}

	return filteredParts, nil
}

// matchesPart проверяет, соответствует ли деталь всем критериям фильтра
// Логика: AND между полями фильтра, OR внутри каждого поля
//
// В отличие от реализации через slices.Contains (O(n²)),
// использует внутренний set на основе map для поиска за O(1),
// что обеспечивает общую сложность O(n + m).
//
// n — количество деталей, m — количество элементов фильтра.
func matchesPart(
	part *model.Part,
	filter *model.PartsFilter,
	uuidSet, nameSet, countrySet, tagSet map[string]struct{},
	categorySet map[model.Category]struct{},
) bool {
	if len(filter.Uuids) > 0 {
		if _, found := uuidSet[part.Uuid]; !found {
			return false
		}
	}

	if len(filter.Names) > 0 {
		if _, found := nameSet[part.Name]; !found {
			return false
		}
	}

	if len(filter.Categories) > 0 {
		if _, found := categorySet[part.Category]; !found {
			return false
		}
	}

	if len(filter.ManufacturerCountries) > 0 {
		if part.Manufacturer == nil {
			return false
		}
		if _, found := countrySet[part.Manufacturer.Country]; !found {
			return false
		}
	}

	if len(filter.Tags) > 0 {
		if !hasAnyTag(part.Tags, tagSet) {
			return false
		}
	}

	return true
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
