package inmemory

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SuiteRepository struct {
	suite.Suite
	ctx  context.Context //nolint:containedctx
	repo *repository
}

func (s *SuiteRepository) SetupTest() {
	s.ctx = context.Background()
	s.repo = NewRepository()
}

func (s *SuiteRepository) TearDownTest() {
	s.repo = nil
}

func TestRepositorySuite(t *testing.T) {
	suite.Run(t, new(SuiteRepository))
}
