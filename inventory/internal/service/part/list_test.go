package part

import (
	"inventory/internal/model"
)

func (s *SuiteService) TestListPartsAllCases() {
	parts := []*model.Part{
		{
			Uuid:         "uuid-1",
			Name:         "Engine",
			Category:     model.CategoryEngine,
			Tags:         []string{"heavy", "metal"},
			Manufacturer: &model.Manufacturer{Name: "SpaceX", Country: "USA"},
		},
		{
			Uuid:         "uuid-2",
			Name:         "Wing",
			Category:     model.CategoryWing,
			Tags:         []string{"light", "composite"},
			Manufacturer: &model.Manufacturer{Name: "Airbus", Country: "France"},
		},
		{
			Uuid:         "uuid-3",
			Name:         "Fuel Pump",
			Category:     model.CategoryFuel,
			Tags:         []string{"liquid"},
			Manufacturer: &model.Manufacturer{Name: "SpaceX", Country: "USA"},
		},
		{
			Uuid:         "uuid-4",
			Name:         "Window",
			Category:     model.CategoryPorthole,
			Tags:         nil,
			Manufacturer: nil,
		},
	}

	s.partRepository.
		On("ListParts", s.ctx).
		Return(parts, nil)

	tests := []struct {
		name     string
		filter   *model.PartsFilter
		expected []string
	}{
		{
			name:     "nil filter returns all",
			filter:   nil,
			expected: []string{"uuid-1", "uuid-2", "uuid-3", "uuid-4"},
		},
		{
			name:     "empty filter returns all",
			filter:   &model.PartsFilter{},
			expected: []string{"uuid-1", "uuid-2", "uuid-3", "uuid-4"},
		},
		{
			name: "filter by uuid",
			filter: &model.PartsFilter{
				Uuids: []string{"uuid-2"},
			},
			expected: []string{"uuid-2"},
		},
		{
			name: "filter by name single",
			filter: &model.PartsFilter{
				Names: []string{"Engine"},
			},
			expected: []string{"uuid-1"},
		},
		{
			name: "filter by name multiple (OR logic)",
			filter: &model.PartsFilter{
				Names: []string{"Engine", "Fuel Pump"},
			},
			expected: []string{"uuid-1", "uuid-3"},
		},
		{
			name: "filter by category single",
			filter: &model.PartsFilter{
				Categories: []model.Category{model.CategoryFuel},
			},
			expected: []string{"uuid-3"},
		},
		{
			name: "filter by category multiple (OR logic)",
			filter: &model.PartsFilter{
				Categories: []model.Category{model.CategoryWing, model.CategoryPorthole},
			},
			expected: []string{"uuid-2", "uuid-4"},
		},
		{
			name: "filter by manufacturer country",
			filter: &model.PartsFilter{
				ManufacturerCountries: []string{"USA"},
			},
			expected: []string{"uuid-1", "uuid-3"},
		},
		{
			name: "filter by manufacturer with nil manufacturers present",
			filter: &model.PartsFilter{
				ManufacturerCountries: []string{"France"},
			},
			expected: []string{"uuid-2"}, // не падает на uuid-4
		},
		{
			name: "filter by tags single",
			filter: &model.PartsFilter{
				Tags: []string{"light"},
			},
			expected: []string{"uuid-2"},
		},
		{
			name: "filter by multiple tags (OR logic)",
			filter: &model.PartsFilter{
				Tags: []string{"metal", "liquid"},
			},
			expected: []string{"uuid-1", "uuid-3"},
		},
		{
			name: "filter with multiple fields (AND logic across fields)",
			filter: &model.PartsFilter{
				Categories:            []model.Category{model.CategoryFuel},
				ManufacturerCountries: []string{"USA"},
			},
			expected: []string{"uuid-3"},
		},
		{
			name: "filter with non-overlapping fields (returns empty)",
			filter: &model.PartsFilter{
				Categories:            []model.Category{model.CategoryFuel},
				ManufacturerCountries: []string{"France"},
			},
			expected: []string{},
		},
		{
			name: "filter with empty values inside fields (ignored safely)",
			filter: &model.PartsFilter{
				Uuids:      []string{},
				Categories: []model.Category{},
				Tags:       []string{},
			},
			expected: []string{"uuid-1", "uuid-2", "uuid-3", "uuid-4"},
		},
	}

	for _, tc := range tests {
		s.Run(tc.name, func() {
			result, err := s.service.ListParts(s.ctx, tc.filter)
			s.Require().NoError(err)
			got := make([]string, 0, len(result))
			for _, p := range result {
				got = append(got, p.Uuid)
			}
			s.ElementsMatch(tc.expected, got)
		})
	}

	s.partRepository.AssertExpectations(s.T())
}

func (s *SuiteService) TestListPartsNoParts() {
	s.partRepository.
		On("ListParts", s.ctx).
		Return([]*model.Part{}, nil).
		Once()

	result, err := s.service.ListParts(s.ctx, nil)
	s.Require().NoError(err)
	s.Empty(result)
	s.partRepository.AssertExpectations(s.T())
}
