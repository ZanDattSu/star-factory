package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/ZanDattSu/star-factory/inventory/internal/repository/mocks"
)

type SuiteService struct {
	suite.Suite

	ctx context.Context //nolint:containedctx

	partRepository *mocks.PartRepository

	service *service
}

func (s *SuiteService) SetupTest() {
	s.ctx = context.Background()

	s.partRepository = mocks.NewPartRepository(s.T())

	s.service = NewService(s.partRepository)
}

func (s *SuiteService) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(SuiteService))
}
