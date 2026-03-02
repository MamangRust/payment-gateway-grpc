package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	userrole_errors "MamangRust/paymentgatewaygrpc/pkg/errors/user_role_errors"
	"context"
)

type userRoleRepository struct {
	db *db.Queries
}

func NewUserRoleRepository(db *db.Queries) *userRoleRepository {
	return &userRoleRepository{
		db: db,
	}
}

func (r *userRoleRepository) AssignRoleToUser(ctx context.Context, req *requests.CreateUserRoleRequest) (*db.UserRole, error) {
	res, err := r.db.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return nil, userrole_errors.ErrAssignRoleToUser
	}

	return res, nil
}

func (r *userRoleRepository) RemoveRoleFromUser(ctx context.Context, req *requests.RemoveUserRoleRequest) error {
	err := r.db.RemoveRoleFromUser(ctx, db.RemoveRoleFromUserParams{
		UserID: int32(req.UserId),
		RoleID: int32(req.RoleId),
	})

	if err != nil {
		return userrole_errors.ErrRemoveRole
	}

	return nil
}
