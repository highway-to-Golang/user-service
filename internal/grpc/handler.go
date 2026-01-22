package grpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/highway-to-Golang/user-service/api/proto/gen/go/user"
	"github.com/highway-to-Golang/user-service/internal/domain"
	apperrors "github.com/highway-to-Golang/user-service/internal/errors"
	"github.com/highway-to-Golang/user-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserHandler struct {
	user.UnimplementedUserServiceServer
	uc *usecase.UseCase
}

func NewUserHandler(uc *usecase.UseCase) *UserHandler {
	return &UserHandler{
		uc: uc,
	}
}

func extractIdempotencyKey(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get("idempotency-key")
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, domain.ErrNotFound) {
		return status.Error(codes.NotFound, "user not found")
	}

	if errors.Is(err, apperrors.ErrRequestAlreadyInProgress) {
		return status.Error(codes.Aborted, "request already in progress")
	}

	slog.Error("unexpected error", "error", err)
	return status.Error(codes.Internal, "internal server error")
}

func (h *UserHandler) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.User, error) {
	idempotencyKey := extractIdempotencyKey(ctx)

	domainReq := domain.CreateUserRequest{
		Name:  req.Name,
		Email: req.Email,
		Role:  req.Role,
	}

	domainUser, err := h.uc.CreateUser(ctx, idempotencyKey, domainReq)
	if err != nil {
		return nil, mapError(err)
	}

	return &user.User{
		Id:        domainUser.ID,
		Name:      domainUser.Name,
		Email:     domainUser.Email,
		Role:      domainUser.Role,
		CreatedAt: timestamppb.New(domainUser.CreatedAt),
		UpdatedAt: timestamppb.New(domainUser.UpdatedAt),
	}, nil
}

func (h *UserHandler) GetUser(ctx context.Context, req *user.GetUserRequest) (*user.User, error) {
	domainUser, err := h.uc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, mapError(err)
	}

	return &user.User{
		Id:        domainUser.ID,
		Name:      domainUser.Name,
		Email:     domainUser.Email,
		Role:      domainUser.Role,
		CreatedAt: timestamppb.New(domainUser.CreatedAt),
		UpdatedAt: timestamppb.New(domainUser.UpdatedAt),
	}, nil
}

func (h *UserHandler) GetAllUsers(ctx context.Context, req *user.GetAllUsersRequest) (*user.GetAllUsersResponse, error) {
	domainUsers, err := h.uc.GetAllUsers(ctx)
	if err != nil {
		return nil, mapError(err)
	}

	protoUsers := make([]*user.User, 0, len(domainUsers))
	for _, u := range domainUsers {
		protoUsers = append(protoUsers, &user.User{
			Id:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: timestamppb.New(u.CreatedAt),
			UpdatedAt: timestamppb.New(u.UpdatedAt),
		})
	}

	return &user.GetAllUsersResponse{
		Users: protoUsers,
		Total: int32(len(domainUsers)),
	}, nil
}

func (h *UserHandler) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.User, error) {
	domainReq := domain.UpdateUserRequest{
		Role: req.Role,
	}
	if req.Name != nil {
		domainReq.Name = req.Name
	}
	if req.Email != nil {
		domainReq.Email = req.Email
	}

	domainUser, err := h.uc.UpdateUser(ctx, req.Id, domainReq)
	if err != nil {
		return nil, mapError(err)
	}

	return &user.User{
		Id:        domainUser.ID,
		Name:      domainUser.Name,
		Email:     domainUser.Email,
		Role:      domainUser.Role,
		CreatedAt: timestamppb.New(domainUser.CreatedAt),
		UpdatedAt: timestamppb.New(domainUser.UpdatedAt),
	}, nil
}

func (h *UserHandler) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	err := h.uc.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, mapError(err)
	}

	return &user.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}
