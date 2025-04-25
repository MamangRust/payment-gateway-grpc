package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type userRoleRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.UserRoleRecordMapping
}

func NewUserRoleRepository(db *db.Queries, ctx context.Context, mapping recordmapper.UserRoleRecordMapping) *userRoleRepository {
	return &userRoleRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *userRoleRepository) AssignRoleToUser(req *requests.CreateUserRoleRequest) (*record.UserRoleRecord, error) {
	res, err := r.db.AssignRoleToUser(r.ctx, db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("cannot assign role: user ID %d or role ID %d not found", req.UserId, req.RoleId)
		}
		return nil, fmt.Errorf("failed to assign role to user ID %d and role ID %d: %w", req.UserId, req.RoleId, err)
	}

	return r.mapping.ToUserRoleRecord(res), nil
}

func (r *userRoleRepository) RemoveRoleFromUser(req *requests.RemoveUserRoleRequest) error {
	err := r.db.RemoveRoleFromUser(r.ctx, db.RemoveRoleFromUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("cannot remove role: no assignment found for user ID %d and role ID %d", req.UserId, req.RoleId)
		}
		return fmt.Errorf("failed to remove role from user ID %d and role ID %d: %w", req.UserId, req.RoleId, err)
	}

	return nil
}
