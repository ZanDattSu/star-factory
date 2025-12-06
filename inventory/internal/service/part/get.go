package part

import (
	"context"

	"go.uber.org/zap"

	"github.com/ZanDattSu/star-factory/inventory/internal/model"
	"github.com/ZanDattSu/star-factory/platform/pkg/logger"
)

func (s *service) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	logger.Debug(ctx, "Getting part",
		zap.String("part_uuid", uuid),
	)

	part, err := s.repository.GetPart(ctx, uuid)
	if err != nil {
		logger.Warn(ctx, "Part not found",
			zap.String("part_uuid", uuid),
			zap.Error(err),
		)
		return nil, &model.PartNotFoundError{PartUUID: uuid}
	}

	logger.Debug(ctx, "Part retrieved successfully",
		zap.String("part_uuid", uuid),
		zap.String("part_name", part.Name),
		zap.String("category", string(part.Category)),
	)

	return part, nil
}
