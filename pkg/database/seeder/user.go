package seeder

import (
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/hash"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type userSeeder struct {
	db     *db.Queries
	ctx    context.Context
	hash   hash.HashPassword
	logger logger.LoggerInterface
}

func NewUserSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface, hash hash.HashPassword) *userSeeder {
	return &userSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
		hash:   hash,
	}
}

func (r *userSeeder) Seed() error {
	for i := 1; i <= 40; i++ {
		email := fmt.Sprintf("user_%s@example.com", uuid.NewString())

		password, err := r.hash.HashPassword("password")

		if err != nil {
			return fmt.Errorf("failed generate password")
		}

		user := db.CreateUserParams{
			Firstname: fmt.Sprintf("User%d", i),
			Lastname:  fmt.Sprintf("Last%d", i),
			Email:     email,
			Password:  password,
		}

		createdUser, err := r.db.CreateUser(r.ctx, user)
		if err != nil {
			r.logger.Error("failed to seed user", zap.Int("user", i), zap.Error(err))
			return fmt.Errorf("failed to seed user %d: %w", i, err)
		}

		if i <= 20 {
			_, err := r.db.TrashUser(r.ctx, createdUser.UserID)
			if err != nil {
				r.logger.Error("failed to trash user", zap.Int("user", i), zap.Error(err))
				return fmt.Errorf("failed to trash user %d: %w", i, err)
			}
		}
	}

	r.logger.Info("user seeded successfully")

	return nil
}
