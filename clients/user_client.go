package clients

import (
	"context"

	"github.com/DarioKnezovic/api-gateway/proto/user"
	"google.golang.org/grpc"
)

type UserClient struct {
	client user.UserServiceClient
}

func NewUserClient(address string) (*UserClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := user.NewUserServiceClient(conn)
	return &UserClient{
		client: client,
	}, nil
}

func (uc *UserClient) CheckUserExistence(userID string) (bool, error) {
	request := &user.UserExistsRequest{
		UserId: userID,
	}

	response, err := uc.client.CheckUserExists(context.Background(), request)
	if err != nil {
		return false, err
	}

	return response.Exists, nil
}
