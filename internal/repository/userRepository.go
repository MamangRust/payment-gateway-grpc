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

type userRepository struct {
	db      *db.Queries
	ctx     context.Context
	mapping recordmapper.UserRecordMapping
}

func NewUserRepository(db *db.Queries, ctx context.Context, mapping recordmapper.UserRecordMapping) *userRepository {
	return &userRepository{
		db:      db,
		ctx:     ctx,
		mapping: mapping,
	}
}

func (r *userRepository) FindAllUsers(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetUsersWithPaginationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetUsersWithPagination(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no users found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}

		return nil, nil, fmt.Errorf("failed to retrieve users: invalid pagination (page %d, size %d) or search criteria '%s'", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordPagination(res), &totalCount, nil
}

func (r *userRepository) FindById(user_id int) (*record.UserRecord, error) {
	res, err := r.db.GetUserByID(r.ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with ID: %d", user_id)
		}
		return nil, fmt.Errorf("failed to find user with ID %d: %w", user_id, err)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) FindByActive(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetActiveUsersWithPaginationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetActiveUsersWithPagination(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no active users found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}

		return nil, nil, fmt.Errorf("failed to find active users: invalid parameters (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordActivePagination(res), &totalCount, nil
}

func (r *userRepository) FindByTrashed(req *requests.FindAllUsers) ([]*record.UserRecord, *int, error) {
	offset := (req.Page - 1) * req.PageSize

	reqDb := db.GetTrashedUsersWithPaginationParams{
		Column1: req.Search,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	}

	res, err := r.db.GetTrashedUsersWithPagination(r.ctx, reqDb)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, fmt.Errorf("no trashed users found matching the criteria (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
		}

		return nil, nil, fmt.Errorf("failed to find trashed users: invalid parameters (page %d, size %d, search '%s')", req.Page, req.PageSize, req.Search)
	}

	var totalCount int
	if len(res) > 0 {
		totalCount = int(res[0].TotalCount)
	} else {
		totalCount = 0
	}

	return r.mapping.ToUsersRecordTrashedPagination(res), &totalCount, nil
}

func (r *userRepository) FindByEmail(email string) (*record.UserRecord, error) {
	res, err := r.db.GetUserByEmail(r.ctx, email)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to find user with email %s: %w", email, err)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) CreateUser(request *requests.CreateUserRequest) (*record.UserRecord, error) {
	req := db.CreateUserParams{
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	user, err := r.db.CreateUser(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to create user: invalid or incomplete user data")
		}
		return nil, fmt.Errorf("failed to create user: invalid or incomplete user data")
	}

	return r.mapping.ToUserRecord(user), nil
}

func (r *userRepository) UpdateUser(request *requests.UpdateUserRequest) (*record.UserRecord, error) {
	req := db.UpdateUserParams{
		UserID:    int32(*request.UserID),
		Firstname: request.FirstName,
		Lastname:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	}

	res, err := r.db.UpdateUser(r.ctx, req)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to update user ID %d: user not found or invalid update data", request.UserID)
		}

		return nil, fmt.Errorf("failed to update user ID %d: user not found or invalid update data", request.UserID)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) TrashedUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.TrashUser(r.ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to move user ID %d to trash: user not found or already trashed", user_id)
		}

		return nil, fmt.Errorf("failed to move user ID %d to trash: user not found or already trashed", user_id)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) RestoreUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.RestoreUser(r.ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to restore user ID %d: user not found in trash", user_id)
		}

		return nil, fmt.Errorf("failed to restore user ID %d: user not found in trash", user_id)
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) DeleteUserPermanent(user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(r.ctx, int32(user_id))

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("failed to permanently delete user ID %d: user not found", user_id)
		}

		return false, fmt.Errorf("failed to permanently delete user ID %d: user not found", user_id)
	}

	return true, nil
}

func (r *userRepository) RestoreAllUser() (bool, error) {
	err := r.db.RestoreAllUsers(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("no trashed users available to restore")
		}

		return false, fmt.Errorf("no trashed users available to restore")
	}

	return true, nil
}

func (r *userRepository) DeleteAllUserPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentUsers(r.ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("cannot permanently delete all users: operation not allowed")
		}

		return false, fmt.Errorf("cannot permanently delete all users: operation not allowed")
	}
	return true, nil
}
