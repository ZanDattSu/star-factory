package part

import (
	"math/rand/v2"
	"time"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
)

// RandomPartsFilter возвращает реалистичный фильтр для тестов ListParts.
func RandomPartsFilter() *model.PartsFilter {
	// создаём несколько фейковых частей, чтобы из них взять данные для фильтра
	parts := []*model.Part{
		RandomPart(),
		RandomPart(),
		RandomPart(),
	}

	// собираем случайные фильтры
	return &model.PartsFilter{
		Uuids: []string{
			parts[0].Uuid,
			parts[1].Uuid,
		},
		Names: []string{
			parts[0].Name,
			parts[2].Name,
		},
		Categories: []model.Category{
			RandomCategory(),
			RandomCategory(),
		},
		ManufacturerCountries: []string{
			parts[0].Manufacturer.Country,
			parts[1].Manufacturer.Country,
		},
		Tags: RandomTags(),
	}
}

func RandomPart() *model.Part {
	return &model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.ProductName(),
		Description:   gofakeit.ProductDescription(),
		Price:         gofakeit.Price(0, 10000),
		StockQuantity: int64(gofakeit.IntN(100)),
		Category:      RandomCategory(),
		Dimensions:    RandomDimensions(),
		Manufacturer:  RandomManufacturer(),
		Tags:          RandomTags(),
		Metadata:      RandomMetadata(),
		CreatedAt:     RandomCreatedAt(),
		UpdatedAt:     time.Time{},
	}
}

func RandomCategory() model.Category {
	categories := []model.Category{
		model.CategoryWing,
		model.CategoryPorthole,
		model.CategoryFuel,
		model.CategoryEngine,
	}
	return categories[rand.IntN(len(categories))] //nolint:gosec
}

func RandomDimensions() *model.Dimensions {
	return &model.Dimensions{
		Length: gofakeit.Float64Range(1, 50),
		Width:  gofakeit.Float64Range(1, 50),
		Height: gofakeit.Float64Range(1, 50),
		Weight: gofakeit.Float64Range(1, 50),
	}
}

func RandomManufacturer() *model.Manufacturer {
	return &model.Manufacturer{
		Name:    gofakeit.Company(),
		Country: gofakeit.Country(),
		Website: gofakeit.DomainName(),
	}
}

func RandomTags() []string {
	return []string{
		gofakeit.Adjective(),
		gofakeit.Noun(),
		gofakeit.Color(),
	}
}

func RandomMetadata() map[string]*model.Value {
	keys := []string{"material", "power", "version", "tested", "priority"}
	meta := make(map[string]*model.Value, len(keys))

	for _, k := range keys {
		//nolint:gosec
		switch rand.IntN(4) {
		case 0:
			meta[k] = model.NewStringValue(gofakeit.Word())
		case 1:
			meta[k] = model.NewInt64Value(int64(gofakeit.Number(1, 100)))
		case 2:
			meta[k] = model.NewFloat64Value(gofakeit.Float64Range(0.1, 999.9))
		default:
			meta[k] = model.NewBoolValue(gofakeit.Bool())
		}
	}
	return meta
}

func RandomCreatedAt() time.Time {
	return gofakeit.DateRange(time.Now().AddDate(-1, 0, 0), time.Now())
}
