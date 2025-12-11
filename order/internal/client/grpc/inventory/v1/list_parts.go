package v1

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ZanDattSu/star-factory/order/internal/client/converter"
	"github.com/ZanDattSu/star-factory/order/internal/model"
	grpcAuth "github.com/ZanDattSu/star-factory/platform/pkg/grpc/interceptor"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
	inventoryV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, partsFilter model.PartsFilter) ([]*model.Part, error) {
	logger.Info(ctx, "Requesting parts from inventory service",
		zap.Int("parts_count", len(partsFilter.Uuids)),
		zap.Strings("part_uuids", partsFilter.Uuids),
	)

	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	parts, err := c.genClient.ListParts(
		ctx,
		&inventoryV1.ListPartsRequest{
			Filter: converter.PartsFilterToProto(partsFilter),
		},
	)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {
			case codes.NotFound:
				logger.Warn(ctx, "Parts not found in inventory",
					zap.Strings("part_uuids", partsFilter.Uuids),
					zap.String("grpc_code", st.Code().String()),
				)
				return nil, fmt.Errorf("inventory: %w", NewPartsNotFoundError(partsFilter.Uuids))
			case codes.Internal:
				logger.Error(ctx, "Inventory service internal error",
					zap.Strings("part_uuids", partsFilter.Uuids),
					zap.String("grpc_code", st.Code().String()),
					zap.Error(err),
				)
				return nil, fmt.Errorf("inventory internal error: %w", err)
			case codes.Unavailable:
				logger.Error(ctx, "Inventory service unavailable",
					zap.Strings("part_uuids", partsFilter.Uuids),
					zap.String("grpc_code", st.Code().String()),
					zap.Error(err),
				)
				return nil, fmt.Errorf("inventory service unavailable: %w", err)
			}
		}

		logger.Error(ctx, "Failed to get parts from inventory",
			zap.Strings("part_uuids", partsFilter.Uuids),
			zap.Error(err),
		)
		return nil, fmt.Errorf("inventory ListParts failed: %w", err)
	}

	logger.Info(ctx, "Successfully received parts from inventory",
		zap.Int("requested_parts", len(partsFilter.Uuids)),
		zap.Int("received_parts", len(parts.Parts)),
	)

	return converter.PartsToModel(parts.Parts), nil
}
