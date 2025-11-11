package converter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"inventory/internal/model"
	"inventory/internal/repository/converter"
)

func TestPartConvertersAllCases(t *testing.T) {
	now := time.Now()

	part := &model.Part{
		Uuid:          "uuid-1",
		Name:          "Engine Core",
		Description:   "Main propulsion system",
		Price:         9999.99,
		StockQuantity: 3,
		Category:      model.CategoryEngine,
		Dimensions: &model.Dimensions{
			Length: 10.5, Width: 5.2, Height: 4.0, Weight: 120.0,
		},
		Manufacturer: &model.Manufacturer{
			Name:    "SpaceX",
			Country: "USA",
			Website: "https://spacex.com",
		},
		Tags:      []string{"heavy", "core"},
		CreatedAt: now,
		UpdatedAt: now.Add(time.Hour),
		Metadata: map[string]*model.Value{
			"string": model.NewStringValue("val"),
			"int":    model.NewInt64Value(42),
			"float":  model.NewFloat64Value(3.14),
			"bool":   model.NewBoolValue(true),
		},
	}

	// === PartToRepoModel + PartToModel roundtrip ===
	repoPart := converter.PartToRepoModel(part)
	require.NotNil(t, repoPart)
	require.Equal(t, part.Manufacturer.Name, repoPart.Manufacturer.Name)
	require.Equal(t, part.Dimensions.Length, repoPart.Dimensions.Length)
	require.Len(t, repoPart.Metadata, 4)

	backToModel := converter.PartToModel(repoPart)
	require.Equal(t, part.Manufacturer.Country, backToModel.Manufacturer.Country)
	require.Equal(t, part.Metadata["int"].Int64Value, backToModel.Metadata["int"].Int64Value)

	// === nil inputs return nil ===
	require.Nil(t, converter.PartToRepoModel(nil))
	require.Nil(t, converter.PartToModel(nil))
	require.Nil(t, converter.DimensionsToRepoModel(nil))
	require.Nil(t, converter.DimensionsToModel(nil))
	require.Nil(t, converter.ManufacturerToRepoModel(nil))
	require.Nil(t, converter.ManufacturerToModel(nil))
	require.Nil(t, converter.MetadataToRepoModel(nil))
	require.Nil(t, converter.MetadataToModel(nil))

	// === Dimensions roundtrip ===
	d := &model.Dimensions{Length: 1, Width: 2, Height: 3, Weight: 4}
	require.Equal(t, d, converter.DimensionsToModel(converter.DimensionsToRepoModel(d)))

	// === Manufacturer roundtrip ===
	m := &model.Manufacturer{Name: "ACME", Country: "DE", Website: "acme.de"}
	require.Equal(t, m, converter.ManufacturerToModel(converter.ManufacturerToRepoModel(m)))

	// === Parts slice conversion ===
	list := []*model.Part{part}
	repoList := converter.PartsToRepoModel(list)
	require.Len(t, repoList, 1)
	modelList := converter.PartsToModel(repoList)
	require.Len(t, modelList, 1)

	// === Metadata: nil и пропуски nil-значений ===
	mdata := map[string]*model.Value{
		"a": model.NewStringValue("x"),
		"b": nil, // пропускается
	}
	out := converter.MetadataToRepoModel(mdata)
	require.Contains(t, out, "a")
	require.NotContains(t, out, "b")
}

func TestPartsFilterConverters_AllCases(t *testing.T) {
	filter := model.PartsFilter{
		Uuids:                 []string{"1", "2"},
		Names:                 []string{"Engine", "Wing"},
		Categories:            []model.Category{model.CategoryEngine, model.CategoryWing},
		ManufacturerCountries: []string{"USA", "FR"},
		Tags:                  []string{"heavy"},
	}

	repoFilter := converter.PartsFilterToRepoModel(filter)
	require.Len(t, repoFilter.Categories, len(filter.Categories))
	require.Equal(t, repoFilter.Uuids, filter.Uuids)
	require.Equal(t, repoFilter.Tags, filter.Tags)

	backToModel := converter.PartsFilterToModel(repoFilter)
	require.Equal(t, filter.Categories, backToModel.Categories)
	require.Equal(t, filter.ManufacturerCountries, backToModel.ManufacturerCountries)
}

func TestCategoryConverters(t *testing.T) {
	cats := []model.Category{model.CategoryEngine, model.CategoryFuel}
	repoCats := converter.CategoriesToRepo(cats)
	require.Len(t, repoCats, 2)

	back := converter.CategoriesFromRepo(repoCats)
	require.Equal(t, cats, back)
}
