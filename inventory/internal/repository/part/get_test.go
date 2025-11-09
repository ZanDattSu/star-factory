package part

import "inventory/internal/model"

func (s *SuiteRepository) TestGetPartSuccess() {
	part, ok := s.repo.GetPart(s.ctx, "uuid-1")
	s.True(ok)
	s.Equal("Engine A", part.Name)
	s.Equal(model.CategoryEngine, part.Category)
}

func (s *SuiteRepository) TestGetPartNotFound() {
	part, ok := s.repo.GetPart(s.ctx, "unknown-uuid")
	s.False(ok)
	s.Nil(part)
}
