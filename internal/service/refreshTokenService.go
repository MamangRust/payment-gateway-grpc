package service

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	"MamangRust/paymentgatewaygrpc/internal/repository"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	refreshtoken_errors "MamangRust/paymentgatewaygrpc/pkg/errors/refresh_token_errors"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"context"

	"go.uber.org/zap"
)

type refreshTokenService struct {
	refreshTokenRepository repository.RefreshTokenRepository
	logger                 logger.LoggerInterface
}

func NewRefreshTokenService(refreshTokenRepository repository.RefreshTokenRepository, logger logger.LoggerInterface) *refreshTokenService {
	return &refreshTokenService{
		refreshTokenRepository: refreshTokenRepository,
		logger:                 logger,
	}
}

func (r *refreshTokenService) FindByToken(ctx context.Context, token string) (*db.RefreshToken, error) {
	refreshToken, err := r.refreshTokenRepository.FindByToken(ctx, token)

	if err != nil {
		r.logger.Error("Failed to find refresh token", zap.Error(err))

		return nil, refreshtoken_errors.ErrFailedFindByToken
	}

	return refreshToken, nil
}

func (r *refreshTokenService) FindByUserId(ctx context.Context, user_id int) (*db.RefreshToken, error) {
	refreshToken, err := r.refreshTokenRepository.FindByUserId(ctx, user_id)

	if err != nil {
		r.logger.Error("Failed to find refresh token", zap.Error(err))
		return nil, refreshtoken_errors.ErrFailedFindByUserID
	}

	return refreshToken, nil
}

func (r *refreshTokenService) UpdateRefreshToken(ctx context.Context, req *requests.UpdateRefreshToken) (*db.RefreshToken, error) {
	refreshToken, err := r.refreshTokenRepository.UpdateRefreshToken(ctx, req)

	if err != nil {
		r.logger.Error("Failed to update refresh token", zap.Error(err))
		return nil, refreshtoken_errors.ErrFailedUpdateRefreshToken
	}

	return refreshToken, nil
}

func (r *refreshTokenService) DeleteRefreshToken(ctx context.Context, token string) error {
	err := r.refreshTokenRepository.DeleteRefreshToken(ctx, token)

	if err != nil {
		r.logger.Error("Failed to delete refresh token", zap.Error(err))
		return refreshtoken_errors.ErrFailedDeleteRefreshToken
	}

	return nil
}
