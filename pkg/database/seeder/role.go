package seeder

import (
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/logger"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type roleSeeder struct {
	db     *db.Queries
	ctx    context.Context
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *db.Queries, ctx context.Context, logger logger.LoggerInterface) *roleSeeder {
	return &roleSeeder{
		db:     db,
		ctx:    ctx,
		logger: logger,
	}
}

func (r *roleSeeder) Seed() error {
	roles := []string{
		"ROLE_ADMIN",
		"ROLE_MERCHANT",
		"ROLE_MANAGER",
		"ROLE_USER",
	}

	for _, roleName := range roles {
		role, err := r.db.CreateRole(r.ctx, roleName)
		if err != nil {
			r.logger.Error("failed to seed role", zap.String("roleName", roleName), zap.Error(err))
			return fmt.Errorf("failed to seed role %s: %w", roleName, err)
		}
		r.logger.Debug("role seeded", zap.String("roleName", role.RoleName))
	}

	r.logger.Debug("all roles seeded successfully", zap.Int("totalRoles", len(roles)))

	return nil
}
