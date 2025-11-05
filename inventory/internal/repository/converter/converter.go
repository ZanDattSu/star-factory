package converter

import (
	"maps"

	"inventory/internal/model"
	repoModel "inventory/internal/repository/model"
)

// === Part ===

func PartToRepoModel(part *model.Part) *repoModel.Part {
	return &repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoModel.Category(part.Category),
		Dimensions:    DimensionsToRepoModel(part.Dimensions),
		Manufacturer:  ManufacturerToRepoModel(part.Manufacturer),
		Tags:          copySlice(part.Tags),
		Metadata:      copyMap(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func PartToModel(p *repoModel.Part) *model.Part {
	return &model.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      model.Category(p.Category),
		Dimensions:    DimensionsToModel(p.Dimensions),
		Manufacturer:  ManufacturerToModel(p.Manufacturer),
		Tags:          copySlice(p.Tags),
		Metadata:      copyMap(p.Metadata),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}

// PartsToModel конвертирует []*repoModel.Part → []*model.Part.
func PartsToModel(repoParts []*repoModel.Part) []*model.Part {
	out := make([]*model.Part, 0, len(repoParts))
	for _, p := range repoParts {
		part := PartToModel(p)
		out = append(out, part)
	}
	return out
}

// PartsToRepoModel конвертирует []*model.Part → []*repoModel.Part.
func PartsToRepoModel(parts []*model.Part) []*repoModel.Part {
	out := make([]*repoModel.Part, 0, len(parts))
	for _, p := range parts {
		part := PartToRepoModel(p)
		out = append(out, part)
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
		Uuids:                 copySlice(f.Uuids),
		Names:                 copySlice(f.Names),
		Categories:            categoriesToRepo(f.Categories),
		ManufacturerCountries: copySlice(f.ManufacturerCountries),
		Tags:                  copySlice(f.Tags),
	}
}

func PartsFilterToModel(f repoModel.PartsFilter) model.PartsFilter {
	return model.PartsFilter{
		Uuids:                 copySlice(f.Uuids),
		Names:                 copySlice(f.Names),
		Categories:            categoriesFromRepo(f.Categories),
		ManufacturerCountries: copySlice(f.ManufacturerCountries),
		Tags:                  copySlice(f.Tags),
	}
}

// === Вспомогательные функции ===

func copySlice[T any](sl []T) []T {
	if sl == nil {
		return nil
	}
	dst := make([]T, len(sl))
	copy(dst, sl)
	return dst
}

func copyMap[T comparable](src map[T]T) map[T]T {
	if src == nil {
		return nil
	}
	dst := make(map[T]T, len(src))
	maps.Copy(src, dst)
	return dst
}

func categoriesToRepo(cats []model.Category) []repoModel.Category {
	out := make([]repoModel.Category, len(cats))
	for i, c := range cats {
		out[i] = repoModel.Category(c)
	}
	return out
}

func categoriesFromRepo(cats []repoModel.Category) []model.Category {
	out := make([]model.Category, len(cats))
	for i, c := range cats {
		out[i] = model.Category(c)
	}
	return out
}
