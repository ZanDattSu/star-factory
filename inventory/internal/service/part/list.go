package part

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

// ListParts возвращает отфильтрованный список деталей.
func (s *service) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	var filterFields []zap.Field
	filterFields = append(filterFields, zap.Bool("filter_empty", filterIsEmpty(filter)))
	if filter != nil {
		filterFields = append(filterFields,
			zap.Int("filter_uuids_count", len(filter.Uuids)),
			zap.Int("filter_names_count", len(filter.Names)),
			zap.Int("filter_categories_count", len(filter.Categories)),
			zap.Int("filter_countries_count", len(filter.ManufacturerCountries)),
			zap.Int("filter_tags_count", len(filter.Tags)),
		)
	} else {
		filterFields = append(filterFields,
			zap.Int("filter_uuids_count", 0),
			zap.Int("filter_names_count", 0),
			zap.Int("filter_categories_count", 0),
			zap.Int("filter_countries_count", 0),
			zap.Int("filter_tags_count", 0),
		)
	}
	logger.Debug(ctx, "Listing parts", filterFields...)

	parts, err := s.repository.ListParts(ctx)
	if err != nil {
		logger.Error(ctx, "Failed to list parts from repository",
			zap.Error(err),
		)
		return []*model.Part{}, fmt.Errorf("error listing parts: %w", err)
	}

	logger.Debug(ctx, "Parts retrieved from repository",
		zap.Int("total_parts", len(parts)),
	)

	if filterIsEmpty(filter) {
		logger.Debug(ctx, "Filter is empty, returning all parts",
			zap.Int("parts_count", len(parts)),
		)
		return parts, nil
	}

	// Создаём set'ы для O(1) проверки
	uuidSet := toSet(filter.Uuids)
	nameSet := toSet(filter.Names)
	countrySet := toSet(filter.ManufacturerCountries)
	tagSet := toSet(filter.Tags)
	categorySet := toSet(filter.Categories)

	filteredParts := make([]*model.Part, 0, len(parts))
	for _, part := range parts {
		if matchesPart(part, filter, uuidSet, nameSet, countrySet, tagSet, categorySet) {
			filteredParts = append(filteredParts, part)
		}
	}

	logger.Info(ctx, "Parts filtered successfully",
		zap.Int("total_parts", len(parts)),
		zap.Int("filtered_parts", len(filteredParts)),
	)

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
		if !hasAnyTag(tagSet, part.Tags) {
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
func hasAnyTag(tagSet map[string]struct{}, partTags []string) bool {
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
