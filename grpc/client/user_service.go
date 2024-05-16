package client

import (
	"context"
	"fmt"
	"order-service/grpc/proto"
	"time"

	"github.com/cenkalti/backoff/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	var user *proto.User

	b := backoff.NewConstantBackOff(1 * time.Second)
	maxRetriesBackOff := backoff.WithMaxRetries(b, 3)

	err := backoff.Retry(func() error {
		fmt.Println("✅✅✅✅✅ UserServiceClient.GetUserByID")
		u, err := c.service.GetUserById(context.Background(), req)
		if st, ok := status.FromError(err); ok && st.Code() == codes.Unavailable {
			return err
		}
		user = u
		return backoff.Permanent(err)
	}, maxRetriesBackOff)

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Close closes the gRPC connection
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
