package part

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
	"inventory/internal/model"
)

type SuiteRepository struct {
	suite.Suite
	ctx  context.Context
	repo *repository
}

func (s *SuiteRepository) SetupTest() {
	s.ctx = context.Background()
	s.repo = NewRepository()

	now := time.Now()
	parts := []*model.Part{
		{
			Uuid:     "uuid-1",
			Name:     "Engine A",
			Category: model.CategoryEngine,
			Tags:     []string{"hot", "metal"},
			Manufacturer: &model.Manufacturer{
				Name:    "ACME",
				Country: "USA",
				Website: "acme.com",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:     "uuid-2",
			Name:     "Wing B",
			Category: model.CategoryWing,
			Tags:     []string{"aero"},
			Manufacturer: &model.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "spacex.com",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			Uuid:     "uuid-3",
			Name:     "Fuel Tank C",
			Category: model.CategoryFuel,
			Tags:     []string{"liquid", "tank"},
			Manufacturer: &model.Manufacturer{
				Name:    "Roscosmos",
				Country: "Russia",
				Website: "roscosmos.ru",
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	for _, p := range parts {
		s.repo.PutPart(p.Uuid, p)
	}
}

func (s *SuiteRepository) TearDownTest() {
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(SuiteRepository))
}
