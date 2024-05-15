package client

import (
	"context"
	"order-service/grpc/proto"

	"google.golang.org/grpc"
)

// UserServiceClient represents the gRPC client for the user service
type UserServiceClient struct {
	conn    *grpc.ClientConn
	service proto.UserServiceClient
}

// NewUserServiceClient creates a new instance of UserServiceClient
func NewUserServiceClient(addr string) (*UserServiceClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	client := proto.NewUserServiceClient(conn)
	return &UserServiceClient{
		conn:    conn,
		service: client,
	}, nil
}

// GetUserByID retrieves user information by ID from the user service
func (c *UserServiceClient) GetUserByID(userID uint64) (*proto.User, error) {
	req := &proto.GetUserByIdRequest{
		UserId: userID,
	}
	user, err := c.service.GetUserById(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// Close closes the gRPC connection
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
