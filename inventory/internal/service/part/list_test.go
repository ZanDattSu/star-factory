package part

import (
	"inventory/internal/model"
)

func (s *SuiteService) TestListPartsSuccess() {
	filter := RandomPartsFilter()
	expectedParts := []*model.Part{
		RandomPart(),
		RandomPart(),
	}

	s.partRepository.On("ListParts", s.ctx, filter).Return(expectedParts)

	parts := s.service.ListParts(s.ctx, filter)

	s.Require().Len(parts, len(expectedParts))
	s.Require().Equal(expectedParts, parts)
}
