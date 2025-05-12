package repository

import (
	"MamangRust/paymentgatewaygrpc/internal/domain/record"
	"MamangRust/paymentgatewaygrpc/internal/domain/requests"
	recordmapper "MamangRust/paymentgatewaygrpc/internal/mapper/record"
	db "MamangRust/paymentgatewaygrpc/pkg/database/schema"
	"MamangRust/paymentgatewaygrpc/pkg/errors/user_errors"
	"context"
	"database/sql"
	"errors"
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
		return nil, nil, user_errors.ErrFindAllUsers
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
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
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
		return nil, nil, user_errors.ErrFindActiveUsers
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
		return nil, nil, user_errors.ErrFindTrashedUsers
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
			return nil, user_errors.ErrUserNotFound
		}

		return nil, user_errors.ErrUserNotFound
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
		return nil, user_errors.ErrCreateUser
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
		return nil, user_errors.ErrUpdateUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) TrashedUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.TrashUser(r.ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrTrashedUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) RestoreUser(user_id int) (*record.UserRecord, error) {
	res, err := r.db.RestoreUser(r.ctx, int32(user_id))

	if err != nil {
		return nil, user_errors.ErrRestoreUser
	}

	return r.mapping.ToUserRecord(res), nil
}

func (r *userRepository) DeleteUserPermanent(user_id int) (bool, error) {
	err := r.db.DeleteUserPermanently(r.ctx, int32(user_id))

	if err != nil {
		return false, user_errors.ErrDeleteUserPermanent
	}

	return true, nil
}

func (r *userRepository) RestoreAllUser() (bool, error) {
	err := r.db.RestoreAllUsers(r.ctx)

	if err != nil {
		return false, user_errors.ErrRestoreAllUsers
	}

	return true, nil
}

func (r *userRepository) DeleteAllUserPermanent() (bool, error) {
	err := r.db.DeleteAllPermanentUsers(r.ctx)

	if err != nil {
		return false, user_errors.ErrDeleteAllUsers
	}
	return true, nil
}
