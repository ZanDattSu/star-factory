package v1

import (
	"context"

	converter "github.com/ZanDattSu/star-factory/notification/internal/converter/api"
	"github.com/ZanDattSu/star-factory/notification/internal/model"
	userV1 "github.com/ZanDattSu/star-factory/shared/pkg/proto/user/v1"
)

type client struct {
	genClient userV1.UserServiceClient
}

func NewClient(genClient userV1.UserServiceClient) *client {
	return &client{genClient: genClient}
}

func (c *client) GetUser(ctx context.Context, userUUID string) (*model.User, error) {
	resp, err := c.genClient.GetUser(ctx, &userV1.GetUserRequest{UserUuid: userUUID})
	if err != nil {
		return nil, err
	}

	return converter.UserFromProto(resp.User), err
}
