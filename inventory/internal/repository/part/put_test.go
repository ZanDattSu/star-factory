package part

import "inventory/internal/model"

func (s *SuiteRepository) TestPutPartOverridesExisting() {
	part := &model.Part{
		Uuid:     "uuid-1",
		Name:     "Engine A v2",
		Category: model.CategoryEngine,
	}
	s.repo.PutPart(part.Uuid, part)

	updated, ok := s.repo.GetPart(s.ctx, "uuid-1")
	s.True(ok)
	s.Equal("Engine A v2", updated.Name)
}
