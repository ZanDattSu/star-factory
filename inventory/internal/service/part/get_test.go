package part

import (
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"inventory/internal/model"
)

func (s *SuiteService) TestGetPartSuccess() {
	expectedPart := RandomPart()

	s.partRepository.
		On("GetPart", s.ctx, expectedPart.Uuid).
		Return(expectedPart, true).
		Once()

	part, err := s.service.GetPart(s.ctx, expectedPart.Uuid)
	s.Require().NoError(err)
	s.Require().Equal(expectedPart, part)
}

func (s *SuiteService) TestGetPartNotFound() {
	uuid := gofakeit.UUID()

	s.partRepository.
		On("GetPart", s.ctx, uuid).
		Return((*model.Part)(nil), false).
		Once()

	part, err := s.service.GetPart(s.ctx, uuid)

	s.Require().Nil(part)
	s.Require().Error(err)

	var notFound *model.PartNotFoundError
	s.Require().ErrorAs(err, &notFound)
	s.Require().Equal(uuid, notFound.PartUUID)
	s.Require().Contains(err.Error(), fmt.Sprintf("part with UUID %q not found", uuid))
}
