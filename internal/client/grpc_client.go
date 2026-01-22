package client

import (
	"context"
	"fmt"

	"github.com/highway-to-Golang/user-service/api/proto/gen/go/user"
	"github.com/highway-to-Golang/user-service/internal/domain"
	apperrors "github.com/highway-to-Golang/user-service/internal/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	client user.UserServiceClient
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &GRPCClient{
		conn:   conn,
		client: user.NewUserServiceClient(conn),
	}, nil
}

func (c *GRPCClient) CreateUser(ctx context.Context, req domain.CreateUserRequest, idempotencyKey string) (domain.User, error) {
	ctx = c.addIdempotencyKey(ctx, idempotencyKey)

	protoReq := &user.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	protoResp, err := c.client.CreateUser(ctx, protoReq)
	if err != nil {
		return domain.User{}, c.mapGRPCError(err)
	}

	return c.protoUserToDomain(protoResp), nil
}

func (c *GRPCClient) GetUser(ctx context.Context, id string) (domain.User, error) {
	protoReq := &user.GetUserRequest{
		Id: id,
	}

	protoResp, err := c.client.GetUser(ctx, protoReq)
	if err != nil {
		return domain.User{}, c.mapGRPCError(err)
	}

	return c.protoUserToDomain(protoResp), nil
}

func (c *GRPCClient) GetAllUsers(ctx context.Context) ([]domain.User, error) {
	protoReq := &user.GetAllUsersRequest{}

	protoResp, err := c.client.GetAllUsers(ctx, protoReq)
	if err != nil {
		return nil, c.mapGRPCError(err)
	}

	domainUsers := make([]domain.User, 0, len(protoResp.Users))
	for _, u := range protoResp.Users {
		domainUsers = append(domainUsers, c.protoUserToDomain(u))
	}

	return domainUsers, nil
}

func (c *GRPCClient) UpdateUser(ctx context.Context, id string, req domain.UpdateUserRequest) (domain.User, error) {
	protoReq := &user.UpdateUserRequest{
		Id:   id,
		Role: req.Role,
	}

	if req.Name != nil {
		protoReq.Name = req.Name
	}
	if req.Email != nil {
		protoReq.Email = req.Email
	}

	protoResp, err := c.client.UpdateUser(ctx, protoReq)
	if err != nil {
		return domain.User{}, c.mapGRPCError(err)
	}

	return c.protoUserToDomain(protoResp), nil
}

func (c *GRPCClient) DeleteUser(ctx context.Context, id string) error {
	protoReq := &user.DeleteUserRequest{
		Id: id,
	}

	_, err := c.client.DeleteUser(ctx, protoReq)
	if err != nil {
		return c.mapGRPCError(err)
	}

	return nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) addIdempotencyKey(ctx context.Context, key string) context.Context {
	if key == "" {
		return ctx
	}

	md := metadata.New(map[string]string{
		"idempotency-key": key,
	})

	return metadata.NewOutgoingContext(ctx, md)
}

func (c *GRPCClient) protoUserToDomain(u *user.User) domain.User {
	return domain.User{
		ID:        u.Id,
		Name:      u.Name,
		Email:     u.Email,
		Role:      u.Role,
		CreatedAt: u.CreatedAt.AsTime(),
		UpdatedAt: u.UpdatedAt.AsTime(),
	}
}

func (c *GRPCClient) mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.NotFound:
		return domain.ErrNotFound
	case codes.Aborted:
		return apperrors.ErrRequestAlreadyInProgress
	default:
		return fmt.Errorf("gRPC error: %s", st.Message())
	}
}
