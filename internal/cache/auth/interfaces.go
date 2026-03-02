package auth_cache

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/response"
	"context"
	"time"
)

type IdentityCache interface {
	SetRefreshToken(ctx context.Context, token string, expiration time.Duration)
	GetRefreshToken(ctx context.Context, token string) (string, bool)
	DeleteRefreshToken(ctx context.Context, token string)
	SetCachedUserInfo(ctx context.Context, user *response.UserResponse, expiration time.Duration)
	GetCachedUserInfo(ctx context.Context, userId string) (*response.UserResponse, bool)
	DeleteCachedUserInfo(ctx context.Context, userId string)
}

type LoginCache interface {
	SetCachedLogin(ctx context.Context, email string, data *response.TokenResponse, expiration time.Duration)
	GetCachedLogin(ctx context.Context, email string) (*response.TokenResponse, bool)
}
