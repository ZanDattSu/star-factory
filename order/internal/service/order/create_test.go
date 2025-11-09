package order

import (
	"slices"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	inventoryV1 "order/internal/client/grpc/inventory/v1"
	"order/internal/model"
)

func (s *SuiteService) TestCreateOrderSuccess() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

	listParts := []*model.Part{
		{Uuid: partUuids[0], Price: gofakeit.Price(50, 150)},
		{Uuid: partUuids[1], Price: gofakeit.Price(50, 150)},
		{Uuid: partUuids[2], Price: gofakeit.Price(50, 150)},
	}

	var expectedTotalPrice float64
	for _, part := range listParts {
		expectedTotalPrice += part.Price
	}

	s.inventoryClient.
		On("ListParts", s.ctx, mock.MatchedBy(func(filter model.PartsFilter) bool {
			return slices.Equal(filter.Uuids, partUuids)
		})).
		Return(listParts, nil).
		Once()

	s.orderRepository.
		On("PutOrder", s.ctx, mock.AnythingOfType("string"), mock.MatchedBy(func(order *model.Order) bool {
			return order.UserUUID == userUUID &&
				slices.Equal(order.PartUuids, partUuids) &&
				order.TotalPrice == expectedTotalPrice &&
				order.Status == model.OrderStatusPendingPayment
		})).
		Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().NoError(err)
	s.Require().NotEmpty(orderUUID)
	s.Require().Equal(expectedTotalPrice, totalPrice)
}

func (s *SuiteService) TestCreateOrderEmptyParts() {
	userUUID := gofakeit.UUID()
	var partUuids []string

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().Error(err)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
	s.Require().Contains(err.Error(), "one or more parts not found")
	s.Require().Contains(err.Error(), "empty parts list")

	s.inventoryClient.AssertNotCalled(s.T(), "ListParts", mock.Anything)
	s.orderRepository.AssertNotCalled(s.T(), "PutOrder", mock.Anything, mock.Anything, mock.Anything)
}

func (s *SuiteService) TestCreateOrderPartsLengthMismatch() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

	parts := []*model.Part{
		{Uuid: partUuids[0], Price: 100.0},
		{Uuid: partUuids[1], Price: 200.0},
	}

	s.inventoryClient.
		On("ListParts", s.ctx, mock.Anything).
		Return(parts, nil).
		Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().Error(err)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
	s.Require().Contains(err.Error(), "one or more parts not found")

	s.orderRepository.AssertNotCalled(s.T(), "PutOrder", mock.Anything, mock.Anything, mock.Anything)
}

func (s *SuiteService) TestCreateOrderInventoryPartsNotFoundError() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

	expectedErr := inventoryV1.NewPartsNotFoundError([]string{partUuids[1]})

	s.inventoryClient.
		On("ListParts", s.ctx, mock.MatchedBy(func(filter model.PartsFilter) bool {
			return slices.Equal(filter.Uuids, partUuids)
		})).
		Return(nil, expectedErr).
		Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().Error(err)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)
	s.Require().Contains(err.Error(), "one or more parts not found")
	s.Require().Contains(err.Error(), partUuids[1])

	s.orderRepository.AssertNotCalled(s.T(), "PutOrder", mock.Anything, mock.Anything, mock.Anything)
}

func (s *SuiteService) TestCreateOrderInventoryGenericError() {
	userUUID := gofakeit.UUID()
	partUuids := []string{gofakeit.UUID(), gofakeit.UUID(), gofakeit.UUID()}

	expectedErr := status.Error(codes.Unavailable, "inventory service unavailable")

	s.inventoryClient.
		On("ListParts", s.ctx, mock.MatchedBy(func(filter model.PartsFilter) bool {
			return slices.Equal(filter.Uuids, partUuids)
		})).
		Return(nil, expectedErr).
		Once()

	orderUUID, totalPrice, err := s.service.CreateOrder(s.ctx, userUUID, partUuids)

	s.Require().Error(err)
	s.Require().Empty(orderUUID)
	s.Require().Zero(totalPrice)

	s.Require().Contains(err.Error(), "inventory")
	s.Require().Contains(err.Error(), "unavailable")

	s.orderRepository.AssertNotCalled(s.T(), "PutOrder", mock.Anything, mock.Anything, mock.Anything)
	s.inventoryClient.AssertNumberOfCalls(s.T(), "ListParts", 1)
}
