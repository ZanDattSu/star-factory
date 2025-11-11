package converter

import (
	"inventory/internal/model"
	repoModel "inventory/internal/repository/model"
)

// === Part ===

func PartToRepoModel(p *model.Part) *repoModel.Part {
	if p == nil {
		return nil
	}

	return &repoModel.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      repoModel.Category(p.Category),
		Dimensions:    DimensionsToRepoModel(p.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      MetadataToRepoModel(p.Metadata),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

func PartToModel(p *repoModel.Part) *model.Part {
	if p == nil {
		return nil
	}

	return &model.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      model.Category(p.Category),
		Dimensions:    DimensionsToModel(p.Dimensions),
		Manufacturer:  ManufacturerToModel(p.Manufacturer),
		Tags:          p.Tags,
		Metadata:      MetadataToModel(p.Metadata),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// PartsToModel конвертирует []*repoModel.Part → []*model.Part.
func PartsToModel(repoParts []*repoModel.Part) []*model.Part {
	out := make([]*model.Part, 0, len(repoParts))
	for _, p := range repoParts {
		out = append(out, PartToModel(p))
	}
	return out
}

// PartsToRepoModel конвертирует []*model.Part → []*repoModel.Part.
func PartsToRepoModel(parts []*model.Part) []*repoModel.Part {
	out := make([]*repoModel.Part, 0, len(parts))
	for _, p := range parts {
		if p == nil {
			out = append(out, nil)
			continue
		}
		out = append(out, PartToRepoModel(p))
	}
	return out
}

// === Dimensions ===

func DimensionsToRepoModel(d *model.Dimensions) *repoModel.Dimensions {
	if d == nil {
		return nil
	}
	return &repoModel.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func DimensionsToModel(d *repoModel.Dimensions) *model.Dimensions {
	if d == nil {
		return nil
	}
	return &model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

// === Manufacturer ===

func ManufacturerToRepoModel(m *model.Manufacturer) *repoModel.Manufacturer {
	if m == nil {
		return nil
	}
	return &repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func ManufacturerToModel(m *repoModel.Manufacturer) *model.Manufacturer {
	if m == nil {
		return nil
	}
	return &model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// === PartsFilter ===

func PartsFilterToRepoModel(f model.PartsFilter) repoModel.PartsFilter {
	return repoModel.PartsFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Categories:            CategoriesToRepo(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

func PartsFilterToModel(f repoModel.PartsFilter) model.PartsFilter {
	return model.PartsFilter{
		Uuids:                 f.Uuids,
		Names:                 f.Names,
		Categories:            CategoriesFromRepo(f.Categories),
		ManufacturerCountries: f.ManufacturerCountries,
		Tags:                  f.Tags,
	}
}

// === Metadata ===

// MetadataToRepoModel конвертирует map[string]*model.Value → map[string]*repoModel.Value.
func MetadataToRepoModel(metadata map[string]*model.Value) map[string]*repoModel.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*repoModel.Value, len(metadata))
	for key, value := range metadata {
		if value == nil {
			continue
		}

		switch {
		case value.StringValue != nil:
			result[key] = repoModel.NewStringValue(*value.StringValue)
		case value.Int64Value != nil:
			result[key] = repoModel.NewInt64Value(*value.Int64Value)
		case value.DoubleValue != nil:
			result[key] = repoModel.NewFloat64Value(*value.DoubleValue)
		case value.BoolValue != nil:
			result[key] = repoModel.NewBoolValue(*value.BoolValue)
		default:
			continue
		}
	}

	return result
}

// MetadataToModel конвертирует map[string]*model.Value в map[string]*inventoryV1.Value
func MetadataToModel(metadata map[string]*repoModel.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*model.Value, len(metadata))
	for key, value := range metadata {
		switch {
		case value.StringValue != nil:
			result[key] = model.NewStringValue(*value.StringValue)
		case value.Int64Value != nil:
			result[key] = model.NewInt64Value(*value.Int64Value)
		case value.DoubleValue != nil:
			result[key] = model.NewFloat64Value(*value.DoubleValue)
		case value.BoolValue != nil:
			result[key] = model.NewBoolValue(*value.BoolValue)
		}
	}

	return result
}

// === Categories ===

func CategoriesFromRepo(cats []repoModel.Category) []model.Category {
	out := make([]model.Category, len(cats))
	for i, c := range cats {
		out[i] = model.Category(c)
	}
	return out
}

func CategoriesToRepo(cats []model.Category) []repoModel.Category {
	out := make([]repoModel.Category, len(cats))
	for i, c := range cats {
		out[i] = repoModel.Category(c)
	}
	return out
}
