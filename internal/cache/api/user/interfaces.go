package user_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
)

type UserQueryCache interface {
	GetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUser, bool)
	SetCachedUsersCache(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUser)

	GetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUserDeleteAt, bool)
	SetCachedUserActiveCache(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUserDeleteAt)

	GetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers) (*response.ApiResponsePaginationUserDeleteAt, bool)
	SetCachedUserTrashedCache(ctx context.Context, req *requests.FindAllUsers, data *response.ApiResponsePaginationUserDeleteAt)

	GetCachedUserCache(ctx context.Context, id int) (*response.ApiResponseUser, bool)
	SetCachedUserCache(ctx context.Context, data *response.ApiResponseUser)
}

type UserCommandCache interface {
	DeleteUserCache(ctx context.Context, id int)
}
