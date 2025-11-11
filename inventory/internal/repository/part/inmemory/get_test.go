package inmemory

import (
	"time"

	"inventory/internal/model"
)

func (s *SuiteRepository) TestGetPartSuccess() {
	now := time.Now()
	part := &model.Part{
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
	}

	err := s.repo.PutPart(s.ctx, part.Uuid, part)
	part, gerr := s.repo.GetPart(s.ctx, "uuid-1")

	s.NoError(gerr)
	s.Require().NoError(err)
	s.Equal("Engine A", part.Name)
	s.Equal(model.CategoryEngine, part.Category)
}

func (s *SuiteRepository) TestGetPartNotFound() {
	part, err := s.repo.GetPart(s.ctx, "unknown-uuid")
	s.Error(err)
	s.Nil(part)
}
