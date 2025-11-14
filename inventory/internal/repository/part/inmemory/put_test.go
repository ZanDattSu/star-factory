package inmemory

import (
	"github.com/ZanDattSu/star-factory/inventory/internal/model"
)

func (s *SuiteRepository) TestPutPartOverridesExisting() {
	part := &model.Part{
		Uuid:     "uuid-1",
		Name:     "Engine A v2",
		Category: model.CategoryEngine,
	}

	err := s.repo.PutPart(s.ctx, part.Uuid, part)
	s.Require().NoError(err)

	updated, err := s.repo.GetPart(s.ctx, "uuid-1")

	s.NoError(err)
	s.Require().NoError(err)
	s.Equal("Engine A v2", updated.Name)
}
