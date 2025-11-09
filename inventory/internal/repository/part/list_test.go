package part

import (
	"inventory/internal/model"
	"inventory/internal/repository/converter"
	repoModel "inventory/internal/repository/model"
)

func (s *SuiteRepository) TestListPartsSuccess() {
	filter := &model.PartsFilter{
		Categories:            []model.Category{model.CategoryEngine, model.CategoryWing},
		ManufacturerCountries: []string{"USA"},
	}

	parts := s.repo.ListParts(s.ctx, filter)

	s.Require().Len(parts, 2)
	s.ElementsMatch(
		[]string{"uuid-1", "uuid-2"},
		[]string{parts[0].Uuid, parts[1].Uuid},
	)
}

func (s *SuiteRepository) TestListPartsFilterByTags() {
	filter := &model.PartsFilter{Tags: []string{"aero"}}
	parts := s.repo.ListParts(s.ctx, filter)

	s.Require().Len(parts, 1)
	s.Equal("uuid-2", parts[0].Uuid)
}

func (s *SuiteRepository) TestListPartsFilterByCountry() {
	filter := &model.PartsFilter{ManufacturerCountries: []string{"Russia"}}
	parts := s.repo.ListParts(s.ctx, filter)

	s.Require().Len(parts, 1)
	s.Equal("uuid-3", parts[0].Uuid)
}

func (s *SuiteRepository) TestListPartsEmptyFilterReturnsAll() {
	result := s.repo.ListParts(s.ctx, &model.PartsFilter{})
	s.Len(result, 3)
}

func (s *SuiteRepository) TestListPartsNilFilterReturnsAll() {
	result := s.repo.ListParts(s.ctx, nil)
	s.Len(result, 3)
}

func (s *SuiteRepository) TestConverterIntegration() {
	p := &model.Part{
		Uuid:     "uuid-999",
		Name:     "Test Part",
		Category: model.CategoryPorthole,
	}

	s.repo.PutPart(p.Uuid, p)
	repoPart, ok := s.repo.parts[p.Uuid]
	s.True(ok)

	s.IsType(&repoModel.Part{}, repoPart)

	modelPart := converter.PartToModel(repoPart)
	s.Equal(p.Uuid, modelPart.Uuid)
	s.Equal(p.Name, modelPart.Name)
}
