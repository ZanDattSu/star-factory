package inmemory

import (
	"github.com/ZanDattSu/star-factory/inventory/internal/model"
	"github.com/ZanDattSu/star-factory/inventory/internal/service/part"
)

func (s *SuiteRepository) TestListPartsSuccess() {
	part1 := part.RandomPart()
	part2 := part.RandomPart()

	expectedParts := []*model.Part{part1, part2}

	err := s.repo.PutPart(s.ctx, expectedParts[0].Uuid, expectedParts[0])
	s.Require().NoError(err)

	err = s.repo.PutPart(s.ctx, expectedParts[1].Uuid, expectedParts[1])
	s.Require().NoError(err)

	parts, err := s.repo.ListParts(s.ctx)

	s.Require().NoError(err)
	s.Require().Len(parts, len(expectedParts))
	s.Require().Equal(expectedParts, parts)
}
