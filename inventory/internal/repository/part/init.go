package part

import (
	"time"

	"github.com/google/uuid"
	repoModel "inventory/internal/repository/model"
)

func (r *repository) InitTestData() {
	parts := []*repoModel.Part{
		{
			Uuid:          uuid.NewString(),
			Name:          "Двигатель",
			Description:   "Мощный ракетный двигатель",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      1, // например, 1 = ENGINE
			Dimensions: &repoModel.Dimensions{
				Length: 200,
				Width:  100,
				Height: 120,
				Weight: 300,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "RocketMotors",
				Country: "Russia",
				Website: "https://rocketmotors.example.com",
			},
			Tags:      []string{"основной", "мотор"},
			Metadata:  map[string]*repoModel.Value{"серия": repoModel.NewStringValue("X100")},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Топливный бак",
			Description:   "Бак для хранения топлива",
			Price:         8000.0,
			StockQuantity: 8,
			Category:      2, // например, 2 = FUEL
			Dimensions: &repoModel.Dimensions{
				Length: 150,
				Width:  150,
				Height: 200,
				Weight: 200,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "FuelTech",
				Country: "Germany",
				Website: "https://fueltech.example.com",
			},
			Tags:      []string{"топливо", "бак"},
			Metadata:  map[string]*repoModel.Value{"материал": repoModel.NewStringValue("титан")},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Иллюминатор",
			Description:   "Прочный иллюминатор для ракеты",
			Price:         3000.0,
			StockQuantity: 15,
			Category:      3, // например, 3 = PORTHOLE
			Dimensions: &repoModel.Dimensions{
				Length: 50,
				Width:  50,
				Height: 10,
				Weight: 20,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "GlassSpace",
				Country: "USA",
				Website: "https://glassspace.example.com",
			},
			Tags:      []string{"стекло", "иллюминатор"},
			Metadata:  map[string]*repoModel.Value{"прозрачность": repoModel.NewFloat64Value(99.9)},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Крыло",
			Description:   "Аэродинамическое крыло",
			Price:         5000.0,
			StockQuantity: 12,
			Category:      4, // например, 4 = WING
			Dimensions: &repoModel.Dimensions{
				Length: 300,
				Width:  50,
				Height: 20,
				Weight: 50,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "WingPro",
				Country: "France",
				Website: "https://wingpro.example.com",
			},
			Tags:      []string{"крыло", "аэродинамика"},
			Metadata:  map[string]*repoModel.Value{"тип": repoModel.NewStringValue("стабилизатор")},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Uuid:          uuid.NewString(),
			Name:          "Панель управления",
			Description:   "Электронная панель управления",
			Price:         7000.0,
			StockQuantity: 7,
			Category:      0, // например, 0 = UNKNOWN
			Dimensions: &repoModel.Dimensions{
				Length: 80,
				Width:  40,
				Height: 10,
				Weight: 10,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "ControlSys",
				Country: "Japan",
				Website: "https://controlsys.example.com",
			},
			Tags:      []string{"электроника", "панель"},
			Metadata:  map[string]*repoModel.Value{"версия": repoModel.NewStringValue("2.1")},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, part := range parts {
		r.parts[part.Uuid] = part
	}
}
