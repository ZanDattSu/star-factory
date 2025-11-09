package converter

import (
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"inventory/internal/model"
)

// === Part ===

// PartToProto конвертирует model.Part в protobuf Part
func PartToProto(part *model.Part) *inventoryV1.Part {
	if part == nil {
		return nil
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToProto(part.Category),
		Dimensions:    DimensionsToProto(part.Dimensions),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     timestamppb.New(part.UpdatedAt),
	}
}

// PartToModel конвертирует protobuf Part в model.Part
func PartToModel(part *inventoryV1.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     part.UpdatedAt.AsTime(),
	}
}

// PartsToProto конвертирует []*model.Part → []*inventoryV1.Part
func PartsToProto(parts []*model.Part) []*inventoryV1.Part {
	out := make([]*inventoryV1.Part, 0, len(parts))
	for _, p := range parts {
		part := PartToProto(p)
		out = append(out, part)
	}
	return out
}

// PartsToModel конвертирует []*inventoryV1.Part → []*model.Part
func PartsToModel(protoParts []*inventoryV1.Part) []*model.Part {
	out := make([]*model.Part, 0, len(protoParts))
	for _, p := range protoParts {
		part := PartToModel(p)
		out = append(out, part)
	}
	return out
}

// === Dimensions ===

// DimensionsToProto конвертирует model.Dimensions в protobuf Dimensions
func DimensionsToProto(dimensions *model.Dimensions) *inventoryV1.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// DimensionsToModel конвертирует protobuf Dimensions в model.Dimensions
func DimensionsToModel(dimensions *inventoryV1.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// === Manufacturer ===

// ManufacturerToProto конвертирует model.Manufacturer в protobuf Manufacturer
func ManufacturerToProto(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// ManufacturerToModel конвертирует protobuf Manufacturer в model.Manufacturer
func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// === Category ===

// CategoryToProto конвертирует model.Category в protobuf Category
func CategoryToProto(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

// CategoryToModel конвертирует protobuf Category в model.Category
func CategoryToModel(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

// === Metadata ===

// MetadataToProto конвертирует map[string]*model.Value в map[string]*inventoryV1.Value
func MetadataToProto(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(metadata))
	for key, value := range metadata {
		switch {
		case value.StringValue != nil:
			result[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_StringValue{StringValue: *value.StringValue},
			}
		case value.Int64Value != nil:
			result[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_Int64Value{Int64Value: *value.Int64Value},
			}
		case value.DoubleValue != nil:
			result[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_DoubleValue{DoubleValue: *value.DoubleValue},
			}
		case value.BoolValue != nil:
			result[key] = &inventoryV1.Value{
				Kind: &inventoryV1.Value_BoolValue{BoolValue: *value.BoolValue},
			}
		}
	}

	return result
}

// MetadataToModel конвертирует map[string]*inventoryV1.Value в map[string]*model.Value
func MetadataToModel(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		if value == nil {
			continue
		}

		switch v := value.Kind.(type) {
		case *inventoryV1.Value_StringValue:
			result[key] = model.NewStringValue(v.StringValue)
		case *inventoryV1.Value_Int64Value:
			result[key] = model.NewInt64Value(v.Int64Value)
		case *inventoryV1.Value_DoubleValue:
			result[key] = model.NewFloat64Value(v.DoubleValue)
		case *inventoryV1.Value_BoolValue:
			result[key] = model.NewBoolValue(v.BoolValue)
		}
	}

	return result
}

// === PartsFilter ===

// PartsFilterToProto конвертирует model.PartsFilter в protobuf PartsFilter
func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	return &inventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categoriesToProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// PartsFilterToModel конвертирует protobuf PartsFilter в model.PartsFilter
func PartsFilterToModel(filter *inventoryV1.PartsFilter) *model.PartsFilter {
	if filter == nil {
		return &model.PartsFilter{}
	}

	return &model.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categoriesFromProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// === Вспомогательные функции ===

func categoriesToProto(cats []model.Category) []inventoryV1.Category {
	out := make([]inventoryV1.Category, len(cats))
	for i, c := range cats {
		out[i] = CategoryToProto(c)
	}
	return out
}

func categoriesFromProto(cats []inventoryV1.Category) []model.Category {
	out := make([]model.Category, len(cats))
	for i, c := range cats {
		out[i] = CategoryToModel(c)
	}
	return out
}
